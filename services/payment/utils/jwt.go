package utils

import (
	"errors"
	"fmt"
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

func CreateRetryTokenOtpRequest(refCode string, recCode string) (string, error) {
	JWT_KEY := config.Cfg.JwtKey
	var secretKey = []byte(JWT_KEY)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"recipient": recCode,
		"reference": refCode,
		"exp":       time.Now().Add(5 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
	})

	signedStr, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedStr, nil
}

type Payload struct {
	Reference string
	Recipient string
}

func VerifyRetryToken(tokenStr string) (*Payload, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.Cfg.JwtKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid or expired token")
	}

	recipient, exists := claims["recipient"].(string)
	if !exists {
		return nil, errors.New("Invalid or expired token")
	}

	reference, exists := claims["reference"].(string)
	if !exists {
		return nil, errors.New("Invalid or expired token")
	}

	payload := Payload{
		Reference: reference,
		Recipient: recipient,
	}

	return &payload, nil

}
