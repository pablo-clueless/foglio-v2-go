package main

import (
	"log"
	"os"

	"foglio/v2/src/config"
	"foglio/v2/src/database"
)

func main() {
	log.Println("Starting database seeder...")

	if err := config.InitializeEnvFile(); err != nil {
		log.Fatal("Failed to initialize env file:", err)
	}
	config.InitializeConfig()

	if err := database.InitializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer func() {
		if err := database.CloseDatabase(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err := database.RunSeeds(); err != nil {
		log.Fatal("Failed to run seeds:", err)
	}

	log.Println("Seeding completed successfully!")
	os.Exit(0)
}
