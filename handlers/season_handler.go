package handlers

import (
	"net/http"

	"f1-statshub.v2/database"
	"github.com/gin-gonic/gin"
)

type Entry struct {
	DriverNumber int
	FullName     string
	Team         string
	Country      string
	Count        int
}

func GetSeasonSummary(c *gin.Context) {
	// 1. Top 3 por victorias reales (última posición registrada de cada carrera)
	rows, err := database.DB.Query(`
		SELECT d.driver_number, d.first_name || ' ' || d.last_name AS full_name, d.team_name, d.country_code, COUNT(*) as wins
		FROM positions p
		JOIN drivers d ON d.driver_number = p.driver_number
		JOIN (
			SELECT driver_number, session_key, MAX(date) AS latest_date
			FROM positions
			WHERE position = 1
			GROUP BY session_key
		) AS last_pos
		ON p.driver_number = last_pos.driver_number AND p.session_key = last_pos.session_key AND p.date = last_pos.latest_date
		GROUP BY d.driver_number
		ORDER BY wins DESC
		LIMIT 3
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var topWins []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.DriverNumber, &e.FullName, &e.Team, &e.Country, &e.Count); err == nil {
			topWins = append(topWins, e)
		}
	}

	// 2. Top 3 por vueltas rápidas
	rows, err = database.DB.Query(`
		SELECT d.driver_number, d.first_name || ' ' || d.last_name AS full_name, d.team_name, d.country_code, COUNT(*) as fastest
		FROM laps l
		JOIN drivers d ON d.driver_number = l.driver_number
		WHERE (l.session_key, l.lap_duration) IN (
			SELECT session_key, MIN(lap_duration)
			FROM laps
			WHERE lap_duration < 9999
			GROUP BY session_key
		)
		GROUP BY d.driver_number
		ORDER BY fastest DESC
		LIMIT 3
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var topFastest []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.DriverNumber, &e.FullName, &e.Team, &e.Country, &e.Count); err == nil {
			topFastest = append(topFastest, e)
		}
	}

	// 3. Top 3 por Pole Positions: primeras posiciones por fecha en cada sesión
	rows, err = database.DB.Query(`
		SELECT d.driver_number, d.first_name || ' ' || d.last_name AS full_name, d.team_name, d.country_code, COUNT(*) as poles
		FROM positions p
		JOIN drivers d ON d.driver_number = p.driver_number
		JOIN (
			SELECT session_key, MIN(date) as first_date
			FROM positions
			WHERE position = 1
			GROUP BY session_key
		) as pole_pos
		ON p.session_key = pole_pos.session_key AND p.date = pole_pos.first_date AND p.position = 1
		GROUP BY d.driver_number
		ORDER BY poles DESC
		LIMIT 3
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var topPoles []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.DriverNumber, &e.FullName, &e.Team, &e.Country, &e.Count); err == nil {
			topPoles = append(topPoles, e)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"season":               2024,
		"top_3_winners":        formatEntries(topWins, "wins"),
		"top_3_fastest_laps":   formatEntries(topFastest, "fastest_laps"),
		"top_3_pole_positions": formatEntries(topPoles, "poles"),
	})
}

func formatEntries(entries []Entry, key string) []gin.H {
	var result []gin.H
	for i, e := range entries {
		item := gin.H{
			"position": i + 1,
			"driver":   e.FullName,
			"team":     e.Team,
			"country":  e.Country,
		}
		item[key] = e.Count
		result = append(result, item)
	}
	return result
}
