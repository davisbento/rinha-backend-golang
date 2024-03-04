package migrations

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

type Migrations struct {
	db *pg.DB
}

func NewMigrations(db *pg.DB) *Migrations {
	return &Migrations{db: db}
}

func (m *Migrations) Run() error {
	// check if table clientes exists
	_, err := m.db.Exec("SELECT 1 from clientes")

	if err != nil {
		fmt.Println("trying create table clientes...")
		// // create table clientes
		_, err = m.db.Exec(`
			CREATE TABLE clientes (
				id SERIAL PRIMARY KEY,
				nome VARCHAR(50),
				limite INT
			);
		`)

		if err != nil {
			return err
		}

		// insert some data
		_, err = m.db.Exec(`
			INSERT INTO
				clientes (nome, limite)
			VALUES
				('o barato sai caro', 1000 * 100),
				('zan corp ltda', 800 * 100),
				('les cruders', 10000 * 100),
				('padaria joia de cocaia', 100000 * 100),
				('kid mais', 5000 * 100);
		`)

		if err != nil {
			return err
		}

		fmt.Println("table clientes created successfully.")
	}

	// check if table extrato exists
	_, err = m.db.Exec("SELECT 1 from extratos")

	if err != nil {
		fmt.Println("trying create table extrato...")

		// create table extrato
		_, err = m.db.Exec(`
			CREATE TABLE extratos (
					id SERIAL PRIMARY KEY,
					cliente_id INT REFERENCES clientes(id),
					valor INT,
					tipo VARCHAR(1),
					descricao VARCHAR(100),
					data TIMESTAMP NOT NULL default NOW()
				);
			`)

		if err != nil {
			return err
		}

		fmt.Println("table extrato created successfully.")
	}

	return nil
}
