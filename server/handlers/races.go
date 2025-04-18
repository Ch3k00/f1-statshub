package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListRaces(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Endpoint en construcción"})
}

func GetRaceDetails(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Endpoint en construcción"})
}
