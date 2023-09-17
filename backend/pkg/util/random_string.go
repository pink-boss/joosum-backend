package util

import (
	"github.com/google/uuid"
	"math/rand"
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
