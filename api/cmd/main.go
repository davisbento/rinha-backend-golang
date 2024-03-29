package main

import (
	"davisbento/rinha-backend-golang/api/config"
	"davisbento/rinha-backend-golang/api/handlers"
	redis "davisbento/rinha-backend-golang/api/infra"
	"davisbento/rinha-backend-golang/api/migrations"
	"davisbento/rinha-backend-golang/api/services"
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.NewConfig()

	// Connect to PostgreSQL database
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.DB_HOSTNAME, cfg.DB_PORT),
		User:     cfg.DB_USER,
		Password: cfg.DB_PASSWORD,
		Database: cfg.DB_NAME,
	})

	defer db.Close()

	// Check connection
	_, err := db.Exec("SELECT 1")
	if err != nil {
		fmt.Printf("Error connecting to database: %s \n", err)
	} else {
		fmt.Println("Connected to database successfully.")
	}

	migrations := migrations.NewMigrations(db)

	err = migrations.Run()

	if err != nil {
		fmt.Printf("Error running migrations: %s \n", err)
	} else {
		fmt.Println("Migrations ran successfully.")
	}

	redis := redis.GetInstance(cfg.REDIS_URL)

	clientService := services.NewClientService(db)
	extratoService := services.NewExtratoService(db)
	extratoHandler := handlers.NewExtratoHandler(extratoService, clientService, redis)

	e := echo.New()

	e.Use(middleware.Recover())
	e.HideBanner = true

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Echo! <3")
	})

	e.GET("/clientes/:id/extrato", extratoHandler.GetExtractHandler())
	e.POST("/clientes/:id/transacoes", extratoHandler.PostExtractHandler())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := "3000"

	e.Start(":" + httpPort)
	fmt.Printf("Server running on port %s \n", httpPort)
}
