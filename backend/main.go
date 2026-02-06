package main

import (
	"bpt-knowledge-center/backend/config"
	"bpt-knowledge-center/backend/routes" // Import the new routes package
	"bpt-knowledge-center/backend/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Config
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// 2. Initialize Services
	config.ConnectDB()
	services.InitStorage()

	// 3. Setup Router
	r := routes.SetupRouter()

	// 4. Run Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
