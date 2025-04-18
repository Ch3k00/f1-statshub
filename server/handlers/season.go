package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSeasonSummary(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Endpoint de resumen de temporada en construcción",
	})
}
