package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"f1-statshub.v2/database"
	"f1-statshub.v2/models"
	"github.com/gin-gonic/gin"
)

func ListDrivers(c *gin.Context) {
	rows, err := database.DB.Query(`
        SELECT driver_number, first_name, last_name, team_name, country_code
        FROM drivers
        ORDER BY driver_number
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		if err := rows.Scan(&d.DriverNumber, &d.FirstName, &d.LastName, &d.TeamName, &d.CountryCode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		drivers = append(drivers, d)
	}

	c.JSON(http.StatusOK, drivers)
}

func GetDriverDetails(c *gin.Context) {
	driverIDStr := c.Param("id")
	driverID, err := strconv.Atoi(driverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT s.session_key, s.circuit_short_name, s.country_name, p.position
		FROM positions p
		JOIN sessions s ON p.session_key = s.session_key
		WHERE p.driver_number = ?
	`, driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var raceResults []models.RaceResult
	wins := 0
	top3 := 0
	maxSpeedGlobal := 0.0

	for rows.Next() {
		var rr models.RaceResult
		var country string

		err := rows.Scan(&rr.SessionKey, &rr.CircuitShortName, &country, &rr.Position)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rr.Race = "GP de " + country

		if rr.Position == 1 {
			wins++
			top3++
		} else if rr.Position <= 3 {
			top3++
		}

		// Verificar vuelta rápida
		var minLap float64
		err = database.DB.QueryRow(`
			SELECT MIN(lap_duration) FROM laps WHERE session_key = ?
		`, rr.SessionKey).Scan(&minLap)

		if err == nil {
			var count int
			err = database.DB.QueryRow(`
				SELECT COUNT(*) FROM laps
				WHERE session_key = ? AND driver_number = ? AND lap_duration = ?
			`, rr.SessionKey, driverID, minLap).Scan(&count)
			if err == nil && count > 0 {
				rr.FastestLap = true
			}
		}

		// Velocidad máxima del piloto
		var speed float64
		err = database.DB.QueryRow(`
			SELECT MAX(st_speed) FROM laps
			WHERE driver_number = ? AND session_key = ?
		`, driverID, rr.SessionKey).Scan(&speed)
		if err == nil {
			rr.MaxSpeed = int(speed)
			if speed > maxSpeedGlobal {
				maxSpeedGlobal = speed
			}
		}

		// Mejor vuelta del piloto
		var bestLap sql.NullFloat64
		err = database.DB.QueryRow(`
			SELECT MIN(lap_duration) FROM laps
			WHERE driver_number = ? AND session_key = ?
		`, driverID, rr.SessionKey).Scan(&bestLap)
		if err == nil && bestLap.Valid {
			rr.BestLapDuration = bestLap.Float64
		}

		raceResults = append(raceResults, rr)
	}

	response := models.DriverDetail{
		DriverID: driverID,
		PerformanceSummary: models.PerformanceSummary{
			Wins:     wins,
			Top3:     top3,
			MaxSpeed: int(maxSpeedGlobal),
		},
		RaceResults: raceResults,
	}

	c.JSON(http.StatusOK, response)
}
