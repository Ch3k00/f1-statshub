package initdata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Position struct {
	DriverNumber int    `json:"driver_number"`
	SessionKey   int    `json:"session_key"`
	Position     int    `json:"position"`
	Date         string `json:"date"`
}

func InitPositions(db *sql.DB) error {
	// Crear tabla si no existe
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS positions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		driver_number INTEGER NOT NULL,
		session_key INTEGER NOT NULL,
		position INTEGER NOT NULL,
		date TEXT NOT NULL,
		FOREIGN KEY (driver_number) REFERENCES drivers(driver_number),
		FOREIGN KEY (session_key) REFERENCES sessions(session_key)
	);`)
	if err != nil {
		return err
	}

	// Obtener session_keys
	rows, err := db.Query(`SELECT session_key FROM sessions`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sessionKeys []int
	for rows.Next() {
		var key int
		rows.Scan(&key)
		sessionKeys = append(sessionKeys, key)
	}

	// Descargar posiciones por carrera
	for _, key := range sessionKeys {
		url := fmt.Sprintf("https://api.openf1.org/v1/position?session_key=%d", key)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error consultando posiciones para sesión %d: %v\n", key, err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var positions []Position
		json.Unmarshal(body, &positions)

		for _, p := range positions {
			if p.DriverNumber == 61 {
				continue // omitir Jack Doohan
			}
			_, err := db.Exec(`
				INSERT INTO positions (driver_number, session_key, position, date)
				VALUES (?, ?, ?, ?)`,
				p.DriverNumber, p.SessionKey, p.Position, p.Date,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error insertando posición: %v\n", err)
			}
		}

		fmt.Printf("✔️ Posiciones insertadas para sesión %d\n", key)
	}

	fmt.Println("✅ Tabla positions poblada con éxito.")
	return nil
}
