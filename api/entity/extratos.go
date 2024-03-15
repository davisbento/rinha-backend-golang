package entity

import "time"

type Extrato struct {
	ID        int
	ClienteID int
	Valor     int
	Tipo      string
	Descricao string
	Data      time.Time
}

type ExtratoBodyDTO struct {
	Tipo      string `json:"tipo"`
	Valor     int    `json:"valor"`
	Descricao string `json:"descricao"`
}

type SaldoDTO struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type ExtratoInsertDTO struct {
	ClienteID int
	Valor     int
	Tipo      string
	Descricao string
}

type TransacaoDTO struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type ResponseExtratoGet struct {
	Saldo             SaldoDTO       `json:"saldo"`
	UltimasTransacoes []TransacaoDTO `json:"ultimas_transacoes"`
}
