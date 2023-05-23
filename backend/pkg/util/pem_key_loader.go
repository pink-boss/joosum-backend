package util

import (
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

// LoadPrivateKey pem key 파일로 부터 private key 를 얻음
func LoadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {

	pemKey, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(pemKey)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
