package util

import (
	"crypto/sha256"
	"encoding/hex"
	"joosum-backend/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Role string

const (
	User  Role = "USER"
	Admin      = "ADMIN"
)

// GenerateNewJWTAccessToken func for generate a new JWT access (private) token
// with user ID and permissions.
func GenerateNewJWTAccessToken(roles []Role, email string) (string, error) {
	// Catch JWT secret key from .env file.
	secret := config.GetEnvConfig("jwt_secret")

	// Create a new JWT access token and claims.
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set public claims:
	claims["email"] = email
	claims["roles"] = roles

	// Generate token.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}

// GenerateNewJWTRefreshToken func for generate a new JWT refresh (public) token.
func GenerateNewJWTRefreshToken() (string, error) {
	// Create a new SHA256 hash.
	hash := sha256.New()

	// Create a new now date and time string with salt.
	refresh := config.GetEnvConfig("jwt_refresh") + time.Now().String()

	// See: https://pkg.go.dev/io#Writer.Write
	_, err := hash.Write([]byte(refresh))
	if err != nil {
		// Return error, it refresh token generation failed.
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
