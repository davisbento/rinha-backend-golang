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

type ResponsePostExtract struct {
	Saldo  int `json:"saldo"`
	Limite int `json:"limite"`
}

func (ch *ExtratoHandler) PostExtractHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")

		body := new(entity.ExtratoBodyDTO)

		if err := c.Bind(body); err != nil {
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: "Invalid payload"})
		}

		idInt, err := strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, struct{ Error string }{Error: "Invalid ID"})
		}

		if idInt > 5 || idInt <= 0 {
			// as we are using a mock database, we can't have more than 5 clients
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Client ID"})
		}

		isPayloadValid := services.ValidateExtratoPayload(body)

		if isPayloadValid != nil {
			return c.JSON(http.StatusUnprocessableEntity, struct{ Error string }{Error: isPayloadValid.Error()})
		}

		client, err := ch.ClientService.GetClientById(idInt)

		if err != nil {
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Client not found"})
		}

		// pode ser um crédito ou um débito
		// no caso, positivo ou negativo
		value := services.GetValue(body.Valor, body.Tipo)

		saldo, err := ch.ExtratoService.GetExtratoSumByClienteId(idInt)

		if err != nil {
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Error getting saldo"})
		}

		// o saldo atual + o valor da transação
		newSaldo := saldo + value

		hasEnoughLimit := services.ClientHasEnoughLimit(client.Limite, newSaldo, body.Tipo)

		if !hasEnoughLimit {
			return c.JSON(http.StatusUnprocessableEntity, struct{ Error string }{Error: "Client has no limit"})
		}

		extratoDTO := entity.ExtratoInsertDTO{
			ClienteID: idInt,
			Valor:     value,
			Tipo:      body.Tipo,
			Descricao: body.Descricao,
		}

		err = ch.ExtratoService.InsertExtrato(extratoDTO)

		if err != nil {
			fmt.Printf("Error inserting extrato: %s \n", err)
			errMessage := fmt.Sprintln(err)
			return c.JSON(http.StatusBadRequest, struct{ Error string }{Error: errMessage})
		}

		return c.JSON(http.StatusOK, ResponsePostExtract{
			Saldo:  newSaldo,
			Limite: client.Limite,
		})
	}
}

func (ch *ExtratoHandler) GetExtractHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, struct{ Error string }{Error: "Invalid ID"})
		}

		if idInt > 5 {
			// as we are using a mock database, we can't have more than 5 clients
			return c.JSON(http.StatusNotFound, struct{ Error string }{Error: "Client not found"})
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
