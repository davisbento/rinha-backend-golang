package handlers

import (
	"davisbento/rinha-backend-golang/entity"
	"davisbento/rinha-backend-golang/services"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ExtratoHandler struct {
	ExtratoService *services.ExtratoService
	ClientService  *services.ClientService
}

func NewExtratoHandler(es *services.ExtratoService, cs *services.ClientService) *ExtratoHandler {
	return &ExtratoHandler{
		ExtratoService: es,
		ClientService:  cs,
	}
}

type ExtratoPayload struct {
	Tipo      string `json:"tipo"`
	Valor     int    `json:"valor"`
	Descricao string `json:"descricao"`
}

type ResponsePostExtract struct {
	Saldo  int `json:"saldo"`
	Limite int `json:"limite"`
}

func (ch *ExtratoHandler) PostExtractHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")

		t := new(ExtratoPayload)

		if err := c.Bind(t); err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Invalid payload"})
		}

		idInt, err := strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Invalid ID"})
		}

		client, err := ch.ClientService.GetClientById(idInt)

		if err != nil {
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Client not found"})
		}

		hasEnoughLimit, err := services.ClientHasEnoughLimit(client.Limite, t.Valor, t.Tipo)

		if err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Error checking limit"})
		}

		if !hasEnoughLimit {
			return c.JSON(http.StatusUnprocessableEntity, struct{ Error string }{Error: "Client has no limit"})
		}

		extratoDTO := services.ExtratoInsertDTO{
			ClienteID: idInt,
			Valor:     t.Valor,
			Tipo:      t.Tipo,
			Descricao: t.Descricao,
		}

		err = ch.ExtratoService.InsertExtrato(extratoDTO)

		if err != nil {
			fmt.Printf("Error inserting extrato: %s \n", err)
			errMessage := fmt.Sprintln(err)
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: errMessage})
		}

		saldo, err := ch.ExtratoService.GetExtratoSumByClienteId(idInt)

		if err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Client not found"})
		}

		return c.JSON(http.StatusOK, ResponsePostExtract{
			Saldo:  saldo,
			Limite: client.Limite,
		})
	}
}

func (ch *ExtratoHandler) GetExtractHandler() func(c echo.Context) error {
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

		saldoTotal, err := ch.ExtratoService.GetExtratoSumByClienteId(idInt)

		if err != nil {
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Error getting saldo"})
		}

		last10Transacoes, err := ch.ExtratoService.GetLast10TransacoesByClienteId(idInt)

		if err != nil {
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Error getting transacoes"})
		}

		return c.JSON(http.StatusOK, entity.ResponseExtratoGet{
			Saldo: entity.SaldoDTO{
				Total:       saldoTotal,
				DataExtrato: time.Now(),
				Limite:      client.Limite,
			},
			UltimasTransacoes: last10Transacoes,
		})

	}

}
