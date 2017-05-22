package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echo_log "github.com/labstack/gommon/log"
	"github.com/lavenderx/squirrel/app"
	"github.com/lavenderx/squirrel/app/crypto"
	"github.com/lavenderx/squirrel/app/log"
	"github.com/lavenderx/squirrel/app/model"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// https://jonathanmh.com/building-a-golang-api-with-echo-and-mysql/
// https://www.netlify.com/blog/2016/10/20/building-a-restful-api-in-go/
// https://xiequan.info/go%E4%B8%8Ejson-web-token/

type JWTClaims struct {
	UserId    int64 `json:"user_id"`
	Username  string `json:"user_name"`
	Cellphone string `json:"cellphone"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

// reference: https://zhuanlan.zhihu.com/p/26300634
type httpError struct {
	code    int
	Key     string `json:"error"`
	Message string `json:"message"`
}

func newHTTPError(code int, key string, msg string) *httpError {
	return &httpError{
		code:    code,
		Key:     key,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *httpError) Error() string {
	return e.Key + ": " + e.Message
}

// httpErrorHandler customize echo's HTTP error handler.
func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "InternalServerError"
		msg  string
	)

	if he, ok := err.(*httpError); ok {
		code = he.code
		key = he.Key
		msg = he.Message
	} else if ehe, ok := err.(*echo.HTTPError); ok {
		code = ehe.Code
		key = http.StatusText(code)
		msg = key
	} else if log.IsDebug() {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				log.Error(err)
			}
		} else {
			err := c.JSON(code, newHTTPError(code, key, msg))
			if err != nil {
				log.Error(err)
			}
		}
	}
}

var (
	assetsHandler http.Handler
	config        *app.Config
	mySQLTemplate *app.MySQLTemplate
)

func init() {
	// Load application config
	config = app.LoadConfig()

	// Load assets and set assetsHandler
	assets := app.Assets()
	assetsHandler = http.FileServer(assets.HTTPBox())

	// Initialize log component
	log.Init()

	// Init MySQL & Redis client
	initMySQL(config)
	initRedis(config)
}

func main() {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	e.Logger.SetLevel(echo_log.OFF)

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderWWWAuthenticate,
			echo.HeaderAuthorization},
		AllowMethods: []string{
			echo.GET,
			echo.PUT,
			echo.POST,
			echo.DELETE},
	}))

	// Api group
	apiGroup := e.Group("/api/v1")
	// Configure middleware with the custom claims type for api group
	apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte("secret"),
	}))
	apiGroup.GET("", api)

	// serves the index.html and other static files from rice
	e.GET("/", staticFilesHandler())
	e.GET("/static/*", staticFilesHandler())

	// Login route
	e.POST("/login", login)

	// Add shutdown hook
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Infof("Server will be closed, which is triggered by %s", sig.String())

			// Close redis client
			log.Info("Closing Redis client")
			if err := app.CloseRedisClient(); err != nil {
				log.Error(err)
			}

			// Close MySQL client
			log.Info("Closing MySQL client")
			if err := app.GetMySQLTemplate().Close(); err != nil {
				log.Error(err)
			}

			log.Infof("Server closed on %s", getLocalIP())
			os.Exit(1)
		}
	}()

	// Start server
	address := fmt.Sprintf(":%v", config.ServerConf.Port)
	log.Infof("Squirrel http server started on [::]%v", address)
	e.Logger.Fatal(e.Start(address))
}

func staticFilesHandler() echo.HandlerFunc {
	return echo.WrapHandler(assetsHandler)
}

// curl -i -w "\n" -H "'Content-Type': 'application/json; charset=UTF-8'" -d "username=test&password=testSecret" http://localhost:7000/login
// {
//   "token": "××××××××××××××××"
// }
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	encryptedPasswd := crypto.EncryptPassword([]byte(password))

	if u := mySQLTemplate.GetByNonEmptyFields(&model.User{
		Username: username,
		Password: encryptedPasswd,
	}); u != nil {
		user := u.(*model.User)
		if username == user.Username && encryptedPasswd == user.Password {
			claims := &JWTClaims{
				user.Id,
				user.Username,
				user.Cellphone,
				user.Email,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
				},
			}

			// Create token with claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Generate encoded token and send it as response.
			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, echo.Map{
				"token": t,
			})
		}
	}

	return echo.ErrUnauthorized
}

func (claims JWTClaims) Valid() error {
	if err := claims.StandardClaims.Valid(); err != nil {
		vErr := err.(*jwt.ValidationError)
		return &httpError{
			code:    http.StatusUnauthorized,
			Key:     "TokenValidError",
			Message: vErr.Inner.Error(),
		}
	}

	if claims.UserId > 0 && claims.Username != "" {
		return nil
	}

	return &httpError{
		code:    http.StatusUnauthorized,
		Key:     "TokenValidError",
		Message: "Must provide user_id and user_name",
	}
}

// curl -i -w "\n" -H "Authorization: Bearer $token" http://localhost:7000/api/v1
func api(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	username := claims.Username
	return c.JSON(http.StatusOK, "Hello, "+username+"!")
}

func initMySQL(config *app.Config) {
	app.ConnectToMySQL(config)

	mySQLTemplate = app.GetMySQLTemplate()
	if err := mySQLTemplate.XormEngine().Sync2(new(model.User)); err != nil {
		panic(err)
	}
}

func initRedis(config *app.Config) {
	app.ConnectToRedis(config)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("Oops: %s", err.Error())
		return ""
	}

	var ipv4Addrs = []string{}

	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				//os.Stdout.WriteString(ipNet.IP.String() + "\n")
				ipv4Addrs = append(ipv4Addrs, ipNet.IP.String())
			}
		}
	}

	if len(ipv4Addrs) == 0 {
		return ""
	}

	return ipv4Addrs[0]
}
