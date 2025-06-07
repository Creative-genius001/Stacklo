package utils

import (
	"time"

	"github.com/Creative-genius001/Stacklo/services/payment/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(role string, id string) (string, error) {

	JWT_KEY := config.Cfg.JwtKey
	var secretKey = []byte(JWT_KEY)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"id":   id,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
