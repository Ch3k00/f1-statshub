package handlers

import (
	"database/sql"
	"f1-statshub/server/database"
	"f1-statshub/server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListDrivers(c *gin.Context) {
	rows, err := database.DB.Query(`
        SELECT driver_number, first_name, last_name, name_acronym, team_name, country_code
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
		if err := rows.Scan(&d.DriverNumber, &d.FirstName, &d.LastName, &d.NameAcronym, &d.TeamName, &d.CountryCode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		drivers = append(drivers, d)
	}

	c.JSON(http.StatusOK, drivers)
}

func GetDriverDetails(c *gin.Context) {
	driverID := c.Param("id")

	// Convertir ID a número
	driverNumber, err := strconv.Atoi(driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de piloto inválido"})
		return
	}

	// Obtener información básica del piloto
	var driver models.Driver
	err = database.DB.QueryRow(`
        SELECT driver_number, first_name, last_name, team_name, country_code
        FROM drivers WHERE driver_number = ?
    `, driverNumber).Scan(
		&driver.DriverNumber,
		&driver.FirstName,
		&driver.LastName,
		&driver.TeamName,
		&driver.CountryCode,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Piloto no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Obtener resultados de carreras para este piloto
	rows, err := database.DB.Query(`
        SELECT s.session_key, s.circuit_short_name, s.country_name, p.position
        FROM positions p
        JOIN sessions s ON p.session_key = s.session_key
        WHERE p.driver_number = ?
        ORDER BY s.date_start
    `, driverNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type RaceResult struct {
		SessionKey       int     `json:"session_key"`
		CircuitShortName string  `json:"circuit_short_name"`
		CountryName      string  `json:"country_name"`
		Position         int     `json:"position"`
		FastestLap       bool    `json:"fastest_lap"`
		MaxSpeed         float64 `json:"max_speed"`
		BestLapDuration  float64 `json:"best_lap_duration"`
	}

	var raceResults []RaceResult
	for rows.Next() {
		var rr RaceResult
		if err := rows.Scan(&rr.SessionKey, &rr.CircuitShortName, &rr.CountryName, &rr.Position); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Verificar si tuvo la vuelta más rápida en esta carrera
		var fastestLap bool
		err = database.DB.QueryRow(`
            SELECT EXISTS (
                SELECT 1 FROM laps 
                WHERE session_key = ? AND driver_number = ?
                ORDER BY lap_duration ASC
                LIMIT 1
            ) AND NOT EXISTS (
                SELECT 1 FROM laps 
                WHERE session_key = ? AND driver_number != ?
                AND lap_duration < (
                    SELECT MIN(lap_duration) FROM laps 
                    WHERE session_key = ? AND driver_number = ?
                )
            )
        `, rr.SessionKey, driverNumber, rr.SessionKey, driverNumber, rr.SessionKey, driverNumber).Scan(&fastestLap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rr.FastestLap = fastestLap

		// Obtener velocidad máxima en esta carrera
		err = database.DB.QueryRow(`
            SELECT MAX(st_speed) FROM laps 
            WHERE session_key = ? AND driver_number = ?
        `, rr.SessionKey, driverNumber).Scan(&rr.MaxSpeed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Obtener mejor tiempo de vuelta en esta carrera
		err = database.DB.QueryRow(`
            SELECT MIN(lap_duration) FROM laps 
            WHERE session_key = ? AND driver_number = ?
        `, rr.SessionKey, driverNumber).Scan(&rr.BestLapDuration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		raceResults = append(raceResults, rr)
	}

	// Calcular resumen del piloto
	var wins, top3Finishes int
	var maxSpeed float64

	// Obtener número de victorias
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM positions 
        WHERE driver_number = ? AND position = 1
    `, driverNumber).Scan(&wins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener número de top 3
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM positions 
        WHERE driver_number = ? AND position <= 3
    `, driverNumber).Scan(&top3Finishes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener velocidad máxima
	err = database.DB.QueryRow(`
        SELECT MAX(st_speed) FROM laps 
        WHERE driver_number = ?
    `, driverNumber).Scan(&maxSpeed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Construir respuesta
	response := gin.H{
		"driver_id":    driverNumber,
		"first_name":   driver.FirstName,
		"last_name":    driver.LastName,
		"team_name":    driver.TeamName,
		"country_code": driver.CountryCode,
		"performance_summary": gin.H{
			"wins":           wins,
			"top_3_finishes": top3Finishes,
			"max_speed":      maxSpeed,
		},
		"race_results": raceResults,
	}

	c.JSON(http.StatusOK, response)
}
