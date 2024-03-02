package handlers

import (
	"davisbento/rinha-backend-golang/entity"
	"davisbento/rinha-backend-golang/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ClienteHandler struct {
	ClientService *services.ClientService
}

func NewClienteHandler(cs *services.ClientService) *ClienteHandler {
	return &ClienteHandler{ClientService: cs}
}

func (ch *ClienteHandler) GetExtractHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Invalid ID"})
		}

		client, err := ch.ClientService.GetClientById(idInt)

		if err != nil {
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Client not found"})
		}

		return c.JSON(http.StatusOK, entity.Cliente{
			ID:     client.ID,
			Nome:   client.Nome,
			Limite: client.Limite,
		})
	}
}
