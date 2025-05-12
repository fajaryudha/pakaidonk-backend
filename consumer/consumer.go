package consumers

import (
	"encoding/json"
	"fmt"
	"log"
	"pakaidonk-backend/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp" // Untuk RabbitMQ
)

func ConsumeMessages(db *gorm.DB, myUrl, queueName string) error {
	// Connect ke RabbitMQ
	conn, err := amqp.Dial(myUrl)
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %v", err)
	}
	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var inquiry models.Inquiry
			err := json.Unmarshal(msg.Body, &inquiry)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			log.Printf("data inquiry %s", msg.Body)

			var existing models.Inquiry
			if err := db.Where("partner_reference_no = ? AND partner_id = ?", inquiry.PartnerReferenceNo, inquiry.PartnerID).First(&existing).Error; err == nil {
				log.Printf("Duplicate partnerReferenceNo: %s", inquiry.PartnerReferenceNo)
				continue
			}
			// Simpan data ke database
			if err := db.Create(&inquiry).Error; err != nil {
				log.Printf("Failed to insert inquiry: %v", err)
				continue
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}
