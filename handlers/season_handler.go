package handlers

import (
	"net/http"

	"f1-statshub.v2/database"
	"f1-statshub.v2/models"
	"github.com/gin-gonic/gin"
)

func GetSeasonSummary(c *gin.Context) {
	season := 2024 // fijo en este caso, pero podrías hacer que sea dinámico

	topWinners := getTopResults(`
		SELECT d.first_name || ' ' || d.last_name, d.team_name, d.country_code, COUNT(*) as wins
		FROM positions p
		JOIN drivers d ON d.driver_number = p.driver_number
		WHERE p.position = 1
		GROUP BY d.driver_number
		ORDER BY wins DESC
		LIMIT 3
	`, "wins")

	topFastest := getTopResults(`
		SELECT d.first_name || ' ' || d.last_name, d.team_name, d.country_code, COUNT(*) as fastest_laps
		FROM laps l
		JOIN drivers d ON d.driver_number = l.driver_number
		WHERE (session_key, lap_duration) IN (
			SELECT session_key, MIN(lap_duration)
			FROM laps
			GROUP BY session_key
		)
		GROUP BY d.driver_number
		ORDER BY fastest_laps DESC
		LIMIT 3
	`, "fastest_laps")

	topPoles := getTopResults(`
		SELECT d.first_name || ' ' || d.last_name, d.team_name, d.country_code, COUNT(*) as poles
		FROM positions p
		JOIN drivers d ON d.driver_number = p.driver_number
		JOIN sessions s ON s.session_key = p.session_key
		WHERE p.position = 1 AND s.session_type = 'Qualifying'
		GROUP BY d.driver_number
		ORDER BY poles DESC
		LIMIT 3
	`, "poles")

	c.JSON(http.StatusOK, models.SeasonSummary{
		Season:            season,
		Top3Winners:       topWinners,
		Top3FastestLaps:   topFastest,
		Top3PolePositions: topPoles,
	})
}

func getTopResults(query, col string) []models.SeasonResultItem {
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []models.SeasonResultItem
	pos := 1
	for rows.Next() {
		var r models.SeasonResultItem
		var count int
		err := rows.Scan(&r.Driver, &r.Team, &r.Country, &count)
		if err != nil {
			continue
		}
		r.Position = pos
		switch col {
		case "wins":
			r.Wins = count
		case "fastest_laps":
			r.FastestLaps = count
		case "poles":
			r.Poles = count
		}
		results = append(results, r)
		pos++
	}
	return results
}
