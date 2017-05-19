package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echo_log "github.com/labstack/gommon/log"
	"github.com/lavenderx/squirrel/app"
	"github.com/lavenderx/squirrel/app/crypto"
	"github.com/lavenderx/squirrel/app/log"
	"github.com/lavenderx/squirrel/app/model"
	"gopkg.in/redis.v5"
	"net/http"
	"strings"
	"time"
)

// https://jonathanmh.com/building-a-golang-api-with-echo-and-mysql/
// https://www.netlify.com/blog/2016/10/20/building-a-restful-api-in-go/

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

	// Start server
	address := fmt.Sprintf(":%v", fmt.Sprint(config.ServerConf.Port))
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
//
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
					ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
				},
			}

			// Create token with claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Generate encoded token and send it as response.
			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return err
			}

			setTokenExpireTime(username, t)

			return c.JSON(http.StatusOK, echo.Map{
				"token": t,
			})
		}
	}

	return echo.ErrUnauthorized
}

func (claims JWTClaims) Valid() error {
	if err := claims.StandardClaims.Valid(); err != nil {
		return err
	}

	if claims.UserId > 0 && claims.Username != "" {
		return nil
	}

	return errors.New("Must provide an user ID")
}

// curl -i -w "\n" -H "Authorization: Bearer $token" http://localhost:7000/api/v1
func api(c echo.Context) error {
	c.Request().Header.Get("Authorization")
	token := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1)

	_, err := app.GetRedisClient().Get(crypto.GetMD5Hash(token)).Result()
	if err == redis.Nil {
		return &httpError{
			Key:     "Unauthorized",
			Message: "Token not exists or expired",
		}
	} else if err != nil {
		log.Error(err)
		return &httpError{
			Key:     "InternalServerError",
			Message: err.Error(),
		}
	} else {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JWTClaims)
		username := claims.Username
		return c.String(http.StatusOK, "Welcome "+username+"!")
	}
}

func setTokenExpireTime(username, token string) {
	client := app.GetRedisClient()
	client.Set(crypto.GetMD5Hash(token), username, 1*time.Hour)
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
