package auth

import (
	"context"
	"fmt"
	"joosum-backend/pkg/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     config.GetEnvConfig("googleClientID"),
		ClientSecret: config.GetEnvConfig("googleClientSecret"),
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
