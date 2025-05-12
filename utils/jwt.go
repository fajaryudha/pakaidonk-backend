package utils

import (
	"errors"
	"log"
	"pakaidonk-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("supersecret") // ganti sesuai kebutuhan

func GenerateJWT(sub string) (string, string, error) {
	now := time.Now()
	exp := now.Add(time.Duration(config.Config.JWT.TTLMinutes) * time.Minute)
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"iss": config.Config.JWT.Issuer,
		"sub": sub,
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"jti": jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secretKey)
	return signed, jti, err
}

func ParserToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi metode signing algoritmanya
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	}, jwt.WithLeeway(2*time.Minute))
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, nil, err
	}

	// Memastikan token valid
	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.New("unable to parse claims")
	}

	return token, claims, nil
}
