package main

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lavenderx/squirrel/boxes"
	"net/http"
)

var (
	images         = boxes.Images()
	cssFiles       = boxes.CSS()
	jsFiles        = boxes.JS()
	templates      = boxes.Templates()
	favicon        = boxes.Favicon()
	faviconHandler http.Handler
	assetsHandler  http.Handler
)

func init() {
	faviconHandler = http.FileServer(favicon.HTTPBox())

	// the file server for rice. "static" is the folder where the files come from.
	assetsHandler = http.FileServer(rice.MustFindBox("static").HTTPBox())
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// serves the index.html and favicon.ico from rice
	e.GET("/", echo.WrapHandler(assetsHandler))
	e.GET("/favicon.ico", echo.WrapHandler(faviconHandler))

	// serves other static files
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetsHandler)))

	// Start server
	e.Logger.Fatal(e.Start(":7000"))
}
