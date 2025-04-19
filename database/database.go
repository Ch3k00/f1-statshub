package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
	}
	fmt.Println("✅ Conexión a la base de datos abierta") // <-- añade esto para debug
}
