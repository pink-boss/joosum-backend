package auth

import (
	"context"
	"fmt"

	goOauth "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// validate google access token
// if valid, return true
// else, return false
func ValidateAccessToken(accessToken string) (bool, error) {
	ctx := context.Background()

	clientID := "YOUR_APP_CLIENT_ID"
	clientSecret := "YOUR_APP_CLIENT_SECRET"

	oauthConfig := &goOauth.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
	}

	token := &goOauth.Token{
		AccessToken: accessToken,
	}

	client := oauthConfig.Client(ctx, token)
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return false, fmt.Errorf("unable to create OAuth2 service: %v", err)
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(accessToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return false, fmt.Errorf("unable to verify access token: %v", err)
	}

	// Check if the token's audience matches your app's client ID.
	if tokenInfo.Audience == clientID {
		return true, nil
	}

	return false, fmt.Errorf("access token is not issued by this app")
}
