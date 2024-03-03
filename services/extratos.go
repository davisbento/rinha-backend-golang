package services

import (
	"davisbento/rinha-backend-golang/entity"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type ExtratoService struct {
	DB *pg.DB
}

func NewExtratoService(db *pg.DB) *ExtratoService {
	return &ExtratoService{DB: db}
}

type ExtratoInsertDTO struct {
	ClienteID int
	Valor     int
	Tipo      string
	Descricao string
}

func (es *ExtratoService) InsertExtrato(payload ExtratoInsertDTO) error {
	// validate payload
	if payload.ClienteID <= 0 {
		return fmt.Errorf("ClienteID is required")
	}

	if payload.Valor <= 0 {
		return fmt.Errorf("valor is required")
	}

	if payload.Tipo == "" {
		return fmt.Errorf("tipo is required")
	}

	if payload.Descricao == "" || len(payload.Descricao) > 10 {
		return fmt.Errorf("descricao is required")
	}

	extrato := &entity.Extrato{
		ClienteID: payload.ClienteID,
		Valor:     payload.Valor,
		Tipo:      payload.Tipo,
		Descricao: payload.Descricao,
	}

	_, err := es.DB.Model(extrato).Insert()

	if err != nil {
		fmt.Printf("Error inserting extrato: %s \n", err)
		return err
	}

	return nil
}

func (es *ExtratoService) GetExtratoSumByClienteId(clienteID int) (int, error) {
	var extratos []entity.Extrato

	err := es.DB.Model(&extratos).Where("cliente_id = ?", clienteID).Select()

	if err != nil {
		fmt.Printf("Error getting extrato: %s \n", err)
		return 0, err
	}

	sum := 0

	for _, extrato := range extratos {
		if extrato.Tipo == "c" {
			sum += extrato.Valor
		} else {
			sum -= extrato.Valor
		}
	}

	return sum, nil
}

func (es *ExtratoService) GetLast10TransacoesByClienteId(clienteID int) ([]entity.TransacaoDTO, error) {
	var extratos []entity.Extrato

	err := es.DB.Model(&extratos).Where("cliente_id = ?", clienteID).Order("id DESC").Limit(10).Select()

	if err != nil {
		fmt.Printf("Error getting extrato: %s \n", err)
		return nil, err
	}

	transacoes := make([]entity.TransacaoDTO, 0)

	for _, extrato := range extratos {
		transacao := entity.TransacaoDTO{
			Valor:       extrato.Valor,
			Tipo:        extrato.Tipo,
			Descricao:   extrato.Descricao,
			RealizadaEm: extrato.Data,
		}

		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}
