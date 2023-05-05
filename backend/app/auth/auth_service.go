package auth

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

func ValidateAccessToken(accessToken string) (bool, error) {
	ctx := context.Background()
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	tokenSource := googleOauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return false, fmt.Errorf("unable to verify access token: %v", err)
	}

	return newToken.Valid(), nil
}
