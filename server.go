package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lavenderx/squirrel/app"
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

var (
	assetsHandler http.Handler
	config        *app.Config
)

func init() {
	assets := app.Assets()
	assetsHandler = http.FileServer(assets.HTTPBox())

	config = app.LoadConfig()
}

func main() {
	e := echo.New()

	lvl, err := app.ParseLevel(config.LoggingConf.Level)
	if err != nil {
		panic(err)
	}
	log.SetLevel(lvl)
	e.Logger.SetLevel(lvl)

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
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
