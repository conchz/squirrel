package main

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// the file server for rice. "static" is the folder where the files come from.
	assetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())
	// serves the index.html from rice
	e.GET("/", echo.WrapHandler(assetHandler))

	// serves other static files
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))

	// Start server
	e.Logger.Fatal(e.Start(":7000"))
}
