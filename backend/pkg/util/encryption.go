package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"joosum-backend/pkg/config"
)

// EncryptString는 AES-GCM을 사용해 문자열을 암호화합니다.
// 암호화 키는 config.yml의 aes_key 항목을 사용합니다.
func EncryptString(plain string) (string, error) {
	key := []byte(config.GetEnvConfig("aes_key"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := aesgcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptString은 EncryptString으로 암호화된 값을 복호화합니다.
func DecryptString(cipherText string) (string, error) {
	key := []byte(config.GetEnvConfig("aes_key"))
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("cipher text too short")
	}
	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plain, err := aesgcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// DecryptIfPossible은 복호화 실패 시 원본 값을 그대로 반환합니다.
func DecryptIfPossible(value string) string {
	plain, err := DecryptString(value)
	if err != nil {
		return value
	}
	return plain
}
