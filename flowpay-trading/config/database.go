package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("trading.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	DB = database
	log.Println("Connected to trading.db successfully!")
}
