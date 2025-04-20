package initdata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Session struct {
	SessionKey       int    `json:"session_key"`
	SessionName      string `json:"session_name"`
	SessionType      string `json:"session_type"`
	Location         string `json:"location"`
	CountryName      string `json:"country_name"`
	Year             int    `json:"year"`
	CircuitShortName string `json:"circuit_short_name"`
	DateStart        string `json:"date_start"`
}

func getRealSessionStart(sessionKey int) (string, error) {
	url := fmt.Sprintf("https://api.openf1.org/v1/laps?session_key=%d", sessionKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var laps []map[string]interface{}
	json.Unmarshal(body, &laps)

	var minDate string
	for _, lap := range laps {
		if date, ok := lap["date_start"].(string); ok {
			if minDate == "" || date < minDate {
				minDate = date
			}
		}
	}

	if minDate == "" {
		return "", fmt.Errorf("fecha real no encontrada")
	}
	return minDate, nil
}

func InitSessions(db *sql.DB) error {
	// Crear tabla si no existe
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		session_key INTEGER PRIMARY KEY,
		session_name TEXT NOT NULL,
		session_type TEXT NOT NULL,
		location TEXT NOT NULL,
		country_name TEXT NOT NULL,
		year INTEGER NOT NULL,
		circuit_short_name TEXT NOT NULL,
		date_start TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	// Descargar sesiones tipo "Race"
	url := "https://api.openf1.org/v1/sessions?session_name=Race&year=2024"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var sessions []Session
	json.Unmarshal(body, &sessions)

	for _, s := range sessions {
		// Obtener fecha real de inicio desde /laps
		realStart, err := getRealSessionStart(s.SessionKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ Usando fecha programada para sesión %d (%s): %v\n", s.SessionKey, s.CircuitShortName, err)
			realStart = s.DateStart
		}

		_, err = db.Exec(`
			INSERT OR REPLACE INTO sessions 
			(session_key, session_name, session_type, location, country_name, year, circuit_short_name, date_start)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			s.SessionKey, s.SessionName, s.SessionType, s.Location,
			s.CountryName, s.Year, s.CircuitShortName, realStart,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error insertando sesión %d: %v\n", s.SessionKey, err)
		} else {
			fmt.Printf("✔️ Sesión %s (%s) insertada con fecha %s\n", s.SessionName, s.CircuitShortName, realStart)
		}
	}

	fmt.Println("✅ Tabla sessions poblada con éxito.")
	return nil
}
