package jwtservice

import (
	"api/internal/config"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	CreateJWT(id int) (string, error)
	CheckJWT(token string) error
	GetIdFromJWT(validToken string) (int, error)
}

type jwtService struct {
	cfg *config.Config
}

func NewJWTService(cfg *config.Config) JWTService {
	return jwtService{
		cfg: cfg,
	}
}

func (s jwtService) CreateJWT(id int) (string, error) {
	idStr := strconv.Itoa(id)
	claims := jwt.RegisteredClaims{
		Issuer:    "DemoBank",
		Subject:   idStr,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingToken, err := token.SignedString(s.cfg.SecretKeyJWT)
	if err != nil {
		return "", err
	}

	return signingToken, nil
}

func (s jwtService) CheckJWT(token string) error {
	parseToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid token")
		}

		return s.cfg.SecretKeyJWT, nil
	})

	if err != nil {
		return fmt.Errorf("invalid token - %w", err)
	}

	if !parseToken.Valid {
		return fmt.Errorf("invalid token - %w", err)
	}

	return nil
}

func (s jwtService) GetIdFromJWT(validToken string) (int, error) {
	parseToken, err := jwt.Parse(validToken, func(t *jwt.Token) (interface{}, error) {
		return s.cfg.SecretKeyJWT, nil
	})

	sub, err := parseToken.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	subInt, err := strconv.Atoi(sub)
	if err != nil {
		return 0, err
	}

	return subInt, nil
}
