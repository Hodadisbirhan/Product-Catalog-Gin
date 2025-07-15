package main

import (
	"catalog-gin/config"
	"catalog-gin/router"
	"catalog-gin/seed"
	"log"
	"os"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to DB and run migrations
	config.ConnectDB()

	// Optional: Seed default roles and permissions
	seed.SeedRolesAndPermissions()

	// Start the server
	r := router.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server is running on port %s\n", port)
	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
