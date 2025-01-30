package util

import (
	"encoding/base64"
	"math/rand"

	"github.com/google/uuid"
)

// generate random string
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func CreateId(domain string) string {
	return domain + "-" + uuid.New().String()
}

// generateRandomState : CSRF 방지용 랜덤 문자열
func GenerateRandomState(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
