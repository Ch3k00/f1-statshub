package main

import (
	"fmt"

	"f1-statshub.v2/database"
	"f1-statshub.v2/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB("proxy.db")

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/corredor", handlers.ListDrivers)
		api.GET("/corredor/detalle/:id", handlers.GetDriverDetails) // 👈 ESTA ES LA LÍNEA IMPORTANTE
		api.GET("/carrera", handlers.ListSessions)
		api.GET("/temporada/resumen", handlers.GetSeasonSummary)
	}

	fmt.Println("🚀 Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
