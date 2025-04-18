package database

import (
	"f1-statshub/server/api"
	"f1-statshub/server/models"
	"fmt"
	"log"
)

func SeedInitialData() error {
	if err := seedDrivers(); err != nil {
		return fmt.Errorf("error seeding drivers: %v", err)
	}

	if err := seedRaces(); err != nil {
		return fmt.Errorf("error seeding races: %v", err)
	}

	if err := seedRaceDetails(); err != nil {
		return fmt.Errorf("error seeding race details: %v", err)
	}

	log.Println("✅ Database seeded successfully")
	return nil
}

func seedDrivers() error {
	// Pilotos de la sesión 9574
	drivers9574, err := api.GetDriversBySession(9574)
	if err != nil {
		return err
	}

	requiredNumbers := map[int]bool{
		1: true, 2: true, 3: true, 4: true, 10: true, 11: true, 14: true, 16: true,
		18: true, 20: true, 22: true, 23: true, 24: true, 27: true, 31: true,
		44: true, 55: true, 63: true, 77: true, 81: true,
	}

	for _, driver := range drivers9574 {
		if requiredNumbers[driver.DriverNumber] {
			if err := insertDriver(driver); err != nil {
				return err
			}
		}
	}

	// Pilotos de la sesión 9636
	drivers9636, err := api.GetDriversBySession(9636)
	if err != nil {
		return err
	}

	for _, driver := range drivers9636 {
		if driver.DriverNumber == 30 || driver.DriverNumber == 50 || driver.DriverNumber == 43 {
			if err := insertDriver(driver); err != nil {
				return err
			}
		}
	}

	return nil
}

func seedRaces() error {
	races, err := api.GetRaces(2024)
	if err != nil {
		return err
	}

	for _, race := range races {
		if err := insertRace(race); err != nil {
			return err
		}
	}

	return nil
}

func seedRaceDetails() error {
	rows, err := DB.Query("SELECT session_key FROM sessions")
	if err != nil {
		return err
	}
	defer rows.Close()

	var sessionKeys []int
	for rows.Next() {
		var key int
		if err := rows.Scan(&key); err != nil {
			return err
		}
		sessionKeys = append(sessionKeys, key)
	}

	for _, sessionKey := range sessionKeys {
		positions, err := api.GetPositions(sessionKey)
		if err != nil {
			return err
		}
		for _, position := range positions {
			if err := insertPosition(position); err != nil {
				return err
			}
		}

		laps, err := api.GetLaps(sessionKey)
		if err != nil {
			return err
		}
		for _, lap := range laps {
			if err := insertLap(lap); err != nil {
				return err
			}
		}
	}

	return nil
}

func insertDriver(driver models.Driver) error {
	stmt, err := DB.Prepare(`
		INSERT OR REPLACE INTO drivers 
		(driver_number, first_name, last_name, name_acronym, team_name, country_code)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		driver.DriverNumber,
		driver.FirstName,
		driver.LastName,
		driver.NameAcronym,
		driver.TeamName,
		driver.CountryCode,
	)
	return err
}

func insertRace(race models.Session) error {
	stmt, err := DB.Prepare(`
		INSERT OR REPLACE INTO sessions 
		(session_key, session_name, session_type, location, country_name, year, circuit_short_name, date_start)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		race.SessionKey,
		race.SessionName,
		race.SessionType,
		race.Location,
		race.CountryName,
		race.Year,
		race.CircuitShortName,
		race.DateStart,
	)
	return err
}

func insertPosition(position models.Position) error {
	stmt, err := DB.Prepare(`
		INSERT OR REPLACE INTO positions 
		(driver_number, session_key, position, date)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		position.DriverNumber,
		position.SessionKey,
		position.Position,
		position.Date,
	)
	return err
}

func insertLap(lap models.Lap) error {
	stmt, err := DB.Prepare(`
		INSERT OR REPLACE INTO laps 
		(driver_number, session_key, lap_number, lap_duration, duration_sector_1, duration_sector_2, duration_sector_3, st_speed, date_start)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		lap.DriverNumber,
		lap.SessionKey,
		lap.LapNumber,
		lap.LapDuration,
		lap.DurationSector1,
		lap.DurationSector2,
		lap.DurationSector3,
		lap.StSpeed,
		lap.DateStart,
	)
	return err
}
