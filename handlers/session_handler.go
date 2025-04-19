package handlers

import (
	"net/http"

	"f1-statshub.v2/database"
	"f1-statshub.v2/models"
	"github.com/gin-gonic/gin"
)

func ListSessions(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT session_key, country_name, date_start, year, circuit_short_name
		FROM sessions
		WHERE session_type = 'Race'
		ORDER BY date_start
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var s models.Session
		err := rows.Scan(&s.SessionKey, &s.CountryName, &s.DateStart, &s.Year, &s.CircuitShortName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sessions = append(sessions, s)
	}

	c.JSON(http.StatusOK, sessions)
}
