package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateToken(userId string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(duration).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func GenerateAccessToken(userId string) (string, error) {
	return GenerateToken(userId, 15*time.Minute)
}

func GenerateRefreshToken(userId string) (string, error) {
	return GenerateToken(userId, 7*24*time.Hour)
}
