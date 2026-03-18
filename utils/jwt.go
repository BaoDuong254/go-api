package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	return secretKey, nil
}

func GenerateJWT(email string, userId int64) (string, error) {
	secretKey, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func ValidateJWT(token string) (int64, error) {
	secretKey, err := getJWTSecret()
	if err != nil {
		return 0, err
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected sign method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, errors.New("Could not parse token")
	}
	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, errors.New("Invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Could not parse claims")
	}
	userId := int64(claims["userId"].(float64))
	return userId, nil
}
