package handlers

import (
	"davisbento/rinha-backend-golang/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ExtratoHandler struct {
	ExtratoService *services.ExtratoService
}

func NewExtratoHandler(es *services.ExtratoService) *ExtratoHandler {
	return &ExtratoHandler{ExtratoService: es}
}

func (ch *ExtratoHandler) PostExtractHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Invalid ID"})
		}

		err = ch.ExtratoService.InsertExtrato(services.ExtratoInsert{
			ClienteID: idInt,
			Valor:     100,
			Tipo:      "c",
			Descricao: "Depósito",
		})

		if err != nil {
			fmt.Printf("Error inserting extrato: %s \n", err)
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Client not found"})
		}

		saldo, err := ch.ExtratoService.GetExtratoSumByClienteId(idInt)

		if err != nil {
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Client not found"})
		}

		fmt.Printf("Saldo: %d \n", saldo)

		type response struct {
			Saldo int `json:"saldo"`
		}

		return c.JSON(http.StatusOK, response{Saldo: saldo})
	}
}
