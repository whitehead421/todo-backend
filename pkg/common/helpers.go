package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ICommon interface {
	HashPassword(password string) string
	CheckPasswordHash(password, hash string) bool
	CreateToken(id uint64) (string, error)
	ValidateToken(tokenString string) (id uint64, err error)
	GenerateUUID() string
}

var secretKey = []byte(GetEnvironmentVariables().JwtSecret)

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(id uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(time.Hour).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (id uint64, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idFloat, ok := claims["id"].(float64)
		if !ok {
			return 0, errors.New("id is not a float64")
		}
		return uint64(idFloat), nil
	}

	return
}

func GenerateUUID() string {
	return uuid.New().String()
}
