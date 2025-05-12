package main

import (
	"fmt"
	"pakaidonk-backend/config"
	"pakaidonk-backend/models"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	// Inisialisasi koneksi ke database
	config.InitDB()

	// Migrasi skema database
	if err := config.DB.AutoMigrate(&models.Inquiry{}); err != nil {
		log.Errorf("Migration failed with error: %v", err)
	}

	fmt.Println("Database migrated successfully!")
}
