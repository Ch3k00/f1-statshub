package initdata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Lap struct {
	DriverNumber    int     `json:"driver_number"`
	SessionKey      int     `json:"session_key"`
	LapNumber       int     `json:"lap_number"`
	LapDuration     float64 `json:"lap_duration"`
	DurationSector1 float64 `json:"duration_sector_1"`
	DurationSector2 float64 `json:"duration_sector_2"`
	DurationSector3 float64 `json:"duration_sector_3"`
	StSpeed         float64 `json:"st_speed"`
	DateStart       string  `json:"date_start"`
}

func InitLaps(db *sql.DB) error {
	// Crear tabla si no existe
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS laps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		driver_number INTEGER NOT NULL,
		session_key INTEGER NOT NULL,
		lap_number INTEGER NOT NULL,
		lap_duration REAL,
		duration_sector_1 REAL,
		duration_sector_2 REAL,
		duration_sector_3 REAL,
		st_speed REAL,
		date_start TEXT,
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

	// Descargar y guardar vueltas
	for _, key := range sessionKeys {
		url := fmt.Sprintf("https://api.openf1.org/v1/laps?session_key=%d", key)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error descargando vueltas para sesión %d: %v\n", key, err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var laps []Lap
		json.Unmarshal(body, &laps)

		for _, lap := range laps {
			if lap.DriverNumber == 61 {
				continue // omitir Jack Doohan
			}

			// Reemplazar valores vacíos o cero por 9999 o -1
			lapDuration := lap.LapDuration
			if lapDuration == 0 {
				lapDuration = 9999.0
			}

			sector1 := lap.DurationSector1
			if sector1 == 0 {
				sector1 = 9999.0
			}

			sector2 := lap.DurationSector2
			if sector2 == 0 {
				sector2 = 9999.0
			}

			sector3 := lap.DurationSector3
			if sector3 == 0 {
				sector3 = 9999.0
			}

			stSpeed := lap.StSpeed
			if stSpeed == 0 {
				stSpeed = -1.0
			}

			_, err := db.Exec(`
				INSERT INTO laps (
					driver_number, session_key, lap_number, lap_duration,
					duration_sector_1, duration_sector_2, duration_sector_3,
					st_speed, date_start
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				lap.DriverNumber, lap.SessionKey, lap.LapNumber, lapDuration,
				sector1, sector2, sector3,
				stSpeed, lap.DateStart,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error insertando vuelta: %v\n", err)
			}
		}

		fmt.Printf("✔️ Vueltas insertadas para sesión %d\n", key)
	}

	fmt.Println("✅ Tabla laps poblada con éxito.")
	return nil
}
