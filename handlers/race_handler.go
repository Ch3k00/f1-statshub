package handlers

import (
	"net/http"
	"strconv"

	"f1-statshub.v2/database"
	"f1-statshub.v2/models"
	"github.com/gin-gonic/gin"
)

func GetRaceDetail(c *gin.Context) {
	idStr := c.Param("id")
	raceID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var detail models.RaceDetail

	// Metadata de la carrera
	err = database.DB.QueryRow(`
		SELECT session_key, country_name, date_start, year, circuit_short_name
		FROM sessions
		WHERE session_key = ?
	`, raceID).Scan(&detail.RaceID, &detail.CountryName, &detail.DateStart, &detail.Year, &detail.CircuitShortName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carrera no encontrada"})
		return
	}

	// Resultados únicos por piloto (sin duplicados)
	rows, err := database.DB.Query(`
		SELECT p.position, d.first_name || ' ' || d.last_name, d.team_name, d.country_code
		FROM positions p
		INNER JOIN (
			SELECT driver_number, MAX(id) AS last_id
			FROM positions
			WHERE session_key = ?
			GROUP BY driver_number
		) latest
		ON p.id = latest.last_id
		JOIN drivers d ON d.driver_number = p.driver_number
		WHERE p.session_key = ?
		ORDER BY p.position ASC
	`, raceID, raceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []models.RaceResultEntry
	for rows.Next() {
		var (
			position int
			r        models.RaceResultEntry
		)
		err := rows.Scan(&position, &r.Driver, &r.Team, &r.Country)
		if err == nil {
			r.Position = strconv.Itoa(position)
			results = append(results, r)
		}
	}

	// Marcar al último con "Ultimo"
	if len(results) > 0 {
		results[len(results)-1].Position = "Ultimo"
	}
	detail.Results = results

	// Vuelta más rápida
	database.DB.QueryRow(`
		SELECT d.first_name || ' ' || d.last_name, l.lap_duration, l.duration_sector_1, l.duration_sector_2, l.duration_sector_3
		FROM laps l
		JOIN drivers d ON d.driver_number = l.driver_number
		WHERE l.session_key = ?
		ORDER BY l.lap_duration ASC
		LIMIT 1
	`, raceID).Scan(
		&detail.FastestLap.Driver,
		&detail.FastestLap.Total,
		&detail.FastestLap.Sector1,
		&detail.FastestLap.Sector2,
		&detail.FastestLap.Sector3,
	)

	// Velocidad máxima
	database.DB.QueryRow(`
		SELECT d.first_name || ' ' || d.last_name, MAX(l.st_speed)
		FROM laps l
		JOIN drivers d ON d.driver_number = l.driver_number
		WHERE l.session_key = ?
	`, raceID).Scan(&detail.MaxSpeed.Driver, &detail.MaxSpeed.SpeedKMH)

	c.JSON(http.StatusOK, detail)
}
