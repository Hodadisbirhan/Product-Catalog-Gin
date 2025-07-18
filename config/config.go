package config

import (
	"catalog-gin/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}

func ConnectDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("❌ DB_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	log.Println("✅ Connected to DB")

	DB = db

	// Enable UUID extension (only needed once per DB)
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	// Auto-migrate models
	err = db.AutoMigrate(
		&model.ContentType{},
		&model.Permission{},
		&model.Role{},
		&model.User{},
	)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Database migrated successfully")
}
