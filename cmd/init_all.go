package main

import (
	"database/sql"
	"fmt"
	"log"

	"f1-statshub.v2/initdata"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./proxy.db")
	if err != nil {
		log.Fatal("❌ Error al abrir la base de datos:", err)
	}
	defer db.Close()

	fmt.Println("🚀 Iniciando carga de datos...")

	if err := initdata.InitDrivers(db); err != nil {
		log.Fatal("Error en InitDrivers:", err)
	}
	if err := initdata.InitSessions(db); err != nil {
		log.Fatal("Error en InitSessions:", err)
	}
	if err := initdata.InitPositions(db); err != nil {
		log.Fatal("Error en InitPositions:", err)
	}
	if err := initdata.InitLaps(db); err != nil {
		log.Fatal("Error en InitLaps:", err)
	}

	fmt.Println("✅ Carga de datos completada.")
}
