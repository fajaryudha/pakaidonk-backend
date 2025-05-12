package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GenerateSignature menggunakan private key untuk menandatangani string
func GenerateSignature(data, privateKeyPath string) (string, error) {
	// Muat private key
	privateKey, err := LoadPrivateKey(privateKeyPath)
	if err != nil {
		return "", err
	}

	// Hash string menggunakan SHA-256
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)

	// Tandatangani data yang di-hash dengan private key menggunakan RSA
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hashed)
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// Encode signature menjadi Base64
	return base64.StdEncoding.EncodeToString(signature), nil
}
