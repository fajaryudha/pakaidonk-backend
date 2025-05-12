package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"pakaidonk-backend/config"

	"pakaidonk-backend/utils"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
)

type InquiryRequest struct {
	PartnerReferenceNo       string `json:"partnerReferenceNo"`
	BeneficiaryAccountNumber string `json:"beneficiaryAccountNumber"`
	Amount                   struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	AdditionalInfo struct {
		BeneficiaryBankCode string `json:"beneficiaryBankCode"`
	} `json:"additionalInfo"`
}

func InquiryHandler(c *gin.Context) {
	var request InquiryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Remove "Bearer "
	tokenString := c.GetHeader("Authorization")[7:] // Remove "Bearer "
	log.Printf(tokenString)

	// Parsing JWT token
	token, claims, err := utils.ParserToken(tokenString)
	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid or expired token"})
		return
	}
	log.Printf("JWT Claims: %v", claims)

	partnerID := c.GetHeader("X-PARTNER-ID")
	if partnerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-PARTNER-ID"})
		return
	}

	var existing config.Inquiry

	errDb := config.DB.Where("partner_reference_no = ? AND partner_id = ?", request.PartnerReferenceNo, partnerID).First(&existing).Error
	if errDb == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate partnerReferenceNo"})
		return
	}

	messageData := map[string]interface{}{
		"partner_reference_no":       request.PartnerReferenceNo,
		"partner_id":                 partnerID,
		"beneficiary_account_number": request.BeneficiaryAccountNumber,
		"amount_value":               request.Amount.Value,
		"amount_currency":            request.Amount.Currency,
		"beneficiary_bank_code":      request.AdditionalInfo.BeneficiaryBankCode,
	}
	log.Printf("message Data %d", messageData)

	msgBytes, _ := json.Marshal(messageData)

	err = config.Publish(config.Config.MessageBroker.RabbitMQ.Queue, msgBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send to message broker: %v", err)})
		return
	}

	minifiedBody, err := utils.MinifyJSON(msgBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to minify request body: %v", err)})
		return
	}
	bodyHash := utils.HashBodySHA256(minifiedBody)

	method := "POST"
	path := "/snap/v1.0/emoney/bank-account-inquiry"
	accessToken := tokenString // This should be extracted from claims or the Authorization header
	timestamp := time.Now().Format(time.RFC3339)
	signString := utils.ComposeSignString(method, path, accessToken, bodyHash, timestamp)
	clientSecret := "keys/private_key.pem" // Path ke private key
	signature, err := utils.GenerateSignature(signString, clientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate signature: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "inquiry accepted", "signature": signature})
}
