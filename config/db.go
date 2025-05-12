package config

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

type Inquiry struct {
	ID                       uint   `gorm:"primaryKey"`
	PartnerReferenceNo       string `gorm:"unique"`
	PartnerID                string
	BeneficiaryAccountNumber string
	AmountValue              float64
	AmountCurrency           string
	BeneficiaryBankCode      string
}

func InitDB() {
	var err error
	DB, err = gorm.Open("mysql", "root:@tcp(localhost:3306)/pakaidonk?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB.AutoMigrate(&Inquiry{})

	log.Println("Database connected")
}
