package handler

import (
	"net/http"
	"pakaidonk-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateSignature(c *gin.Context) {
	clientKey := c.GetHeader("X-CLIENT-KEY")
	timestamp := c.GetHeader("X-TIMESTAMP")

	if clientKey == "" || timestamp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing headers"})
		return
	}

	// Validasi format timestamp
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	message := clientKey + "|" + timestamp

	privateKey, err := utils.LoadPrivateKey("keys/private_key.pem")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load private key"})
		return
	}

	signature, err := utils.SignMessage(privateKey, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"signature": signature})
}
