package handler

import (
	"net/http"
	"pakaidonk-backend/config"
	"pakaidonk-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessTokenRequest struct {
	GrantType string `json:"grantType" binding:"required,max=18"`
}

func GenerateAccessToken(c *gin.Context) {
	var req AccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.GrantType != "client_credentials" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grantType must be 'client_credentials'"})
		return
	}

	signature := c.GetHeader("X-SIGNATURE")
	clientKey := c.GetHeader("X-CLIENT-KEY")
	timestamp := c.GetHeader("X-TIMESTAMP")

	if signature == "" || clientKey == "" || timestamp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing headers"})
		return
	}

	message := clientKey + "|" + timestamp

	// Load public key
	pubKey, err := utils.LoadPublicKey("keys/public_key.pem")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid public keys"})
		return
	}

	// Verify signature
	valid := utils.VerifySignature(pubKey, message, signature)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Generate JWT
	tokenStr, jti, err := utils.GenerateJWT(clientKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save to Redis
	err = utils.SaveTokenToRedis(jti, tokenStr, time.Duration(config.Config.JWT.TTLMinutes)*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}
