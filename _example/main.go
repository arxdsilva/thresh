package main

import (
	"github.com/arxdsilva/thresh/_example/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/hc", handlers.HealthCheck)
	e.Logger.Fatal(e.Start(":3000"))
}
