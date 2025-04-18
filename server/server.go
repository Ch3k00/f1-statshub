package main

import (
	"f1-statshub/server/database"
	"f1-statshub/server/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("📦 Iniciando servidor F1 StatsHub...")

	// Inicializar base de datos
	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Error initializing database: %v", err)
	}
	log.Println("✅ Base de datos inicializada.")

	// Poblar datos iniciales
	if err := database.SeedInitialData(); err != nil {
		log.Fatalf("❌ Error seeding initial data: %v", err)
	}
	log.Println("✅ Datos cargados exitosamente.")

	// Configurar router
	r := gin.Default()

	// Configurar endpoints
	r.GET("/api/corredor", handlers.ListDrivers)
	r.GET("/api/corredor/detalle/:id", handlers.GetDriverDetails)
	r.GET("/api/carrera", handlers.ListRaces)
	r.GET("/api/carrera/detalle/:id", handlers.GetRaceDetails)
	r.GET("/api/temporada/resumen", handlers.GetSeasonSummary)

	// Confirmar inicio
	log.Println("🚀 Servidor corriendo en http://localhost:8080")

	// Iniciar servidor
	// Iniciar servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
