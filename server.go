package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echo_log "github.com/labstack/gommon/log"
	"github.com/lavenderx/squirrel/app"
	"github.com/lavenderx/squirrel/app/log"
	"net/http"
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
	initMySQLConnection(config)
	initRedisConnection(config)
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

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// TODO Integrate with Database
	if username == "jon" && password == "shhh!" {
		claims := &JWTClaims{
			1,
			"Jon Snow",
			"15612345678",
			"dolphineor@gmail.com",
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
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
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

func api(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	username := claims.Username
	return c.String(http.StatusOK, "Welcome "+username+"!")
}

func initMySQLConnection(config *app.Config) {
	err := app.ConnectToMySQL(config)
	if err != nil {
		panic(err)
	}
}

func initRedisConnection(config *app.Config) {
	app.ConnectToRedis(config)
}
