package initdata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Driver struct {
	DriverNumber int    `json:"driver_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NameAcronym  string `json:"name_acronym"`
	TeamName     string `json:"team_name"`
	CountryCode  string `json:"country_code"`
}

func InitDrivers(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS drivers (
		driver_number INTEGER PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		name_acronym TEXT NOT NULL,
		team_name TEXT NOT NULL,
		country_code TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	sources := map[string][]int{
		"9574": {1, 2, 3, 4, 10, 11, 14, 16, 18, 20, 22, 23, 24, 27, 31, 44, 55, 63, 77, 81},
		"9636": {30, 43, 50},
	}

	for sessionKey, desiredDrivers := range sources {
		url := fmt.Sprintf("https://api.openf1.org/v1/drivers?session_key=%s", sessionKey)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var allDrivers []Driver
		json.Unmarshal(body, &allDrivers)

		allowed := make(map[int]bool)
		for _, num := range desiredDrivers {
			allowed[num] = true
		}

		for _, d := range allDrivers {
			if !allowed[d.DriverNumber] {
				continue
			}

			_, err := db.Exec(`
				INSERT OR REPLACE INTO drivers (driver_number, first_name, last_name, name_acronym, team_name, country_code)
				VALUES (?, ?, ?, ?, ?, ?)`,
				d.DriverNumber, d.FirstName, d.LastName, d.NameAcronym, d.TeamName, d.CountryCode,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error insertando piloto %d: %v\n", d.DriverNumber, err)
			} else {
				fmt.Printf("✔️ Piloto %s %s insertado\n", d.FirstName, d.LastName)
			}
		}
	}

	fmt.Println("✅ Tabla drivers poblada con éxito.")
	return nil
}
