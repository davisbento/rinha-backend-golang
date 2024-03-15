package services

import (
	"davisbento/rinha-backend-golang/api/entity"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type ExtratoService struct {
	db *pg.DB
}

func NewExtratoService(db *pg.DB) *ExtratoService {
	return &ExtratoService{db: db}
}

func ValidateExtratoPayload(payload *entity.ExtratoBodyDTO) error {
	if payload.Valor <= 0 {
		return fmt.Errorf("valor is required")
	}

	if payload.Tipo == "" {
		return fmt.Errorf("tipo is required")
	}

	if payload.Descricao == "" || len(payload.Descricao) > 10 {
		return fmt.Errorf("descricao is required")
	}

	return nil
}

func (es *ExtratoService) InsertExtrato(payload entity.ExtratoInsertDTO) error {
	extrato := &entity.Extrato{
		ClienteID: payload.ClienteID,
		Valor:     payload.Valor,
		Tipo:      payload.Tipo,
		Descricao: payload.Descricao,
	}

	_, err := es.db.Model(extrato).Insert()

	if err != nil {
		fmt.Printf("Error inserting extrato: %s \n", err)
		return err
	}

	return nil
}

func (es *ExtratoService) GetExtratoSumByClienteId(clienteID int) (int, error) {
	var extratos []entity.Extrato

	err := es.db.Model(&extratos).Where("cliente_id = ?", clienteID).Select()

	if err != nil {
		fmt.Printf("Error getting extrato: %s \n", err)
		return 0, err
	}

	sum := 0

	for _, extrato := range extratos {
		sum += extrato.Valor
	}

	return sum, nil
}

func (es *ExtratoService) GetLast10TransacoesByClienteId(clienteID int) ([]entity.TransacaoDTO, error) {
	var extratos []entity.Extrato

	err := es.db.Model(&extratos).Where("cliente_id = ?", clienteID).Order("id DESC").Limit(10).Select()

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
