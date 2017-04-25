package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lavenderx/squirrel/boxes"
	"net/http"
)

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

	// serves the index.html and favicon.ico from rice
	e.GET("/", echo.WrapHandler(assetsHandler))
	e.GET("/favicon.ico", echo.WrapHandler(assetsHandler))

	// serves other static files
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetsHandler)))

	// Start server
	e.Logger.Fatal(e.Start(":7000"))
}
