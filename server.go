package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lavenderx/squirrel/boxes"
	"net/http"
	"os"
)

// https://jonathanmh.com/building-a-golang-api-with-echo-and-mysql/

var (
	env           boxes.Environment
	assetsHandler http.Handler
)

func init() {
	args := os.Args[1:]
	env = boxes.Environment(args[0])

	assets := boxes.Assets()
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

	// serves the index.html and other static files from rice
	e.GET("/", staticFilesHandler())
	e.GET("/static/*", staticFilesHandler())

	// Start server
	e.Logger.Fatal(e.Start(":7000"))
}

func staticFilesHandler() echo.HandlerFunc {
	return echo.WrapHandler(assetsHandler)
}
