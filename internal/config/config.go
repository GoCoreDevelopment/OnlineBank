package config

import (
	"errors"
	"os"
)

type Config struct {
	Port string
	DataBaseURL string
	SecretKeyJWT []byte
}

func NewConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("satabaseURL it cannot be empty")
	}

	secretKey := os.Getenv("SECRET_KEY_JWT")
	if secretKey == "" {
		return nil, errors.New("secret key it cannot be empty")
	}

	secretKeyByte := []byte(secretKey)

	return &Config{
		Port: port,
		DataBaseURL: databaseURL,
		SecretKeyJWT: secretKeyByte,
	}, nil
}