package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lavenderx/squirrel/boxes"
	"net/http"
)

// https://jonathanmh.com/building-a-golang-api-with-echo-and-mysql/

var (
	assets        = boxes.Assets()
	assetsHandler http.Handler
)

func init() {
	// the file server for rice. "static" is the folder where the files come from.
	assetsHandler = http.FileServer(assets.HTTPBox())
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

	// serves the index.html from rice
	e.GET("/", echo.WrapHandler(assetsHandler))

	// serves other static files
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetsHandler)))

	// Start server
	e.Logger.Fatal(e.Start(":7000"))
}
