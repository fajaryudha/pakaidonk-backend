package main

import (
	"log"
	"pakaidonk-backend/config"
	consumers "pakaidonk-backend/consumer"
	"pakaidonk-backend/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadConfig("config/config.yaml")
	config.InitDB()
	log.Printf("test %s", config.Config.MessageBroker.RabbitMQ.URL)
	err := config.InitRabbitMQ(config.Config.MessageBroker.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	go consumers.ConsumeMessages(
		config.DB,
		config.Config.MessageBroker.RabbitMQ.URL,
		config.Config.MessageBroker.RabbitMQ.Queue,
	)

	r := gin.Default()

	r.POST("/asymmetric-signature", handler.GenerateSignature)
	r.POST("/access-token", handler.GenerateAccessToken)
	r.POST("/snap/v1.0/emoney/bank-account-inquiry", handler.InquiryHandler)

	r.Run(":6000") // listen and serve
}
