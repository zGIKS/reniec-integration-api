// @title ACME Backend API
// @version 1.0
// @description API REST para gestión de citas, catálogo de servicios y clientes con integración RENIEC
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

package main

import (
	"log"
	"net/http"

	"acme/config"
	"acme/database"
	_ "acme/docs"
	"acme/catalog"
	"acme/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (try multiple paths)
	envPaths := []string{".env", "../../../.env", "./.env"}
	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Println("Warning: .env file not found or unable to load, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Set Gin mode based on configuration
	gin.SetMode(cfg.Server.Mode)

	// Connect to database
	db, err := database.Connect(cfg.Database.ConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize database schema
	log.Println("Creating database tables...")
	if err := database.CreateTables(db); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	log.Println("Seeding initial data...")
	if err := database.SeedData(db); err != nil {
		log.Fatal("Failed to seed data:", err)
	}

	// Create services using factory pattern
	serviceFactory := catalog.NewServiceFactory(db, cfg)
	services, err := serviceFactory.CreateServices()
	if err != nil {
		log.Fatal("Failed to create services:", err)
	}

	// Create handlers
	handlers := serviceFactory.CreateHandlers(services)

	// Setup router
	r := gin.Default()
	router.SetupRoutes(r, handlers, cfg)

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port

	if cfg.IsProduction() {
		log.Printf("Production server starting on %s", serverAddr)
		log.Printf("Environment: %s", cfg.App.Environment)
	} else {
		log.Printf("Development server starting on %s", serverAddr)
		log.Printf("Swagger documentation: http://localhost:%s/swagger/index.html", cfg.Server.Port)
		log.Printf("Environment: %s | Gin Mode: %s", cfg.App.Environment, cfg.Server.Mode)
	}

	log.Fatal(http.ListenAndServe(serverAddr, r))
}
