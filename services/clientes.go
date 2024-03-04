package services

import (
	"davisbento/rinha-backend-golang/entity"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type ClientService struct {
	DB *pg.DB
}

func NewClientService(db *pg.DB) *ClientService {
	return &ClientService{DB: db}
}

func (cs *ClientService) GetClientById(id int) (*entity.Cliente, error) {
	client := &entity.Cliente{ID: id}

	err := cs.DB.Model(client).WherePK().Select()

	if err != nil {
		fmt.Printf("Error getting client: %s \n", err)
		return nil, err
	}

	return client, nil
}

func ClientHasEnoughLimit(clientLimit int, valor int, tipo string) (bool, error) {
	if tipo == "c" {
		return true, nil
	}

	// o cliente pode gastar até o dobro do limite em débito
	// ex: se o limite é 100, o cliente pode gastar até 200
	// o limite nao pode ser negativo
	maxLimit := clientLimit * 2

	if maxLimit < valor {
		return false, nil
	}

	return true, nil
}
