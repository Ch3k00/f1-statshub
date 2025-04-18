package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	// Eliminar la BD existente para pruebas (opcional)
	os.Remove("./proxy.db")

	var err error
	DB, err = sql.Open("sqlite3", "./proxy.db")
	if err != nil {
		return err
	}

	// Crear tablas
	if err := createTables(); err != nil {
		return err
	}

	return nil
}

func createTables() error {
	schema, err := os.ReadFile("./server/database/schema.sql")
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(schema))
	return err
}
