package utils

import (
	"os"
	"time"

	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func CreateToken(id string, isVerified bool) (string, error) {

	JWT_KEY := os.Getenv("JWT_KEY")
	var secretKey = []byte(JWT_KEY)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"isVerified": isVerified,
		"exp":        time.Now().Add(72 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		logger.Logger.Warn("jwt token could not be created", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}
