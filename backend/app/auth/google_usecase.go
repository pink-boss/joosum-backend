package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	localConfig "joosum-backend/pkg/config"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleUsecae struct {
}

func (GoogleUsecae) ValidateIdToken(idToken string) (bool, error) {
	ctx := context.Background()

	oauth2Service, err := oauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
	if err != nil {
		return false, fmt.Errorf("unable to create OAuth2 service: %v", err)
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return false, fmt.Errorf("unable to verify id token: %v", err)
	}

	// Check if the token's audience matches your app's client ID.
	if tokenInfo.Audience == localConfig.GetEnvConfig("googleClientID") {
		return true, nil
	}

	return false, fmt.Errorf("id token is not issued by this app")
}

func (GoogleUsecae) ValidateIdTokenForAndroid(idToken string) (bool, error) {
	ctx := context.Background()

	oauth2Service, err := oauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
	if err != nil {
		return false, fmt.Errorf("unable to create OAuth2 service: %v", err)
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return false, fmt.Errorf("unable to verify id token: %v", err)
	}

	// Check if the token's audience matches your app's client ID.
	if tokenInfo.Audience == localConfig.GetEnvConfig("googleAndroidClientID") {
		return true, nil
	}

	return false, fmt.Errorf("id token is not issued by this app")
}

func (GoogleUsecae) GetUserEmail(idToken string) (string, error) {
	ctx := context.Background()

	oauth2Service, err := oauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
	if err != nil {
		return "", fmt.Errorf("unable to create OAuth2 service: %v", err)
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return "", fmt.Errorf("unable to verify id token: %v", err)
	}

	// Check if the token's audience matches your app's client ID.
	if tokenInfo.Audience != localConfig.GetEnvConfig("googleClientID") {
		return "", fmt.Errorf("id token is not issued by this app")
	}

	// Return the user's email address.
	if tokenInfo.Email != "" {
		return tokenInfo.Email, nil
	}

	return "", fmt.Errorf("unable to retrieve user's email")
}

func (GoogleUsecae) GetUserInfoFromToken(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	// 구글 userinfo 엔드포인트
	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo endpoint returned status %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo: %v", err)
	}

	return &userInfo, nil
}
