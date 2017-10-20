package main

import (
	"net/http"

	"github.com/labstack/echo"
	"google.golang.org/appengine"
)

var e *echo.Echo

func main() {
	e = echo.New()

	e.GET("/", serveRoot)

	// pass off app engine to echo
	http.Handle("/", e)
	appengine.Main()
}

func serveRoot(c echo.Context) error {
	return c.String(http.StatusOK, "The coolest 200 OK you've ever seen!")
}
