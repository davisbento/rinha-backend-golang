package main

import (
	"davisbento/rinha-backend-golang/config"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.NewConfig()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Echo! <3")
	})

	e.GET("/clientes/:id/extrato", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, struct{ ID string }{ID: id})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := "3000"

	e.Logger.Fatal(e.Start(":" + httpPort))
}
