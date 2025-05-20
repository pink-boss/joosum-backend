package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	localConfig "joosum-backend/pkg/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleUsecae struct {
}

func (GoogleUsecae) ValidateIdToken(idToken string) (bool, error) {
	ctx := context.Background()

	oauth2Service, err := gOauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
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

	oauth2Service, err := gOauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
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

func (GoogleUsecae) ValidateIdTokenForWeb(idToken string) (bool, error) {
	ctx := context.Background()
	
	// 인증 코드로 보이는 경우 (4/로 시작하는 경우)
	if len(idToken) > 2 && idToken[:2] == "4/" {
		// OAuth2 설정
		config := &oauth2.Config{
			ClientID:     localConfig.GetEnvConfig("googleWebClientID"),
			ClientSecret: localConfig.GetEnvConfig("googleWebSecret"),
			RedirectURL:  localConfig.GetEnvConfig("googleWebRedirect"),
			Endpoint:     google.Endpoint,
		}

		// 인증 코드를 토큰으로 교환
		token, err := config.Exchange(ctx, idToken)
		if err != nil {
			return false, fmt.Errorf("failed to exchange auth code for token: %v", err)
		}

		// 액세스 토큰으로 사용자 정보 가져오기
		userInfoURL := "https://www.googleapis.com/oauth2/v3/userinfo"
		req, err := http.NewRequest("GET", userInfoURL, nil)
		if err != nil {
			return false, fmt.Errorf("failed to create userinfo request: %v", err)
		}
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false, fmt.Errorf("failed to get userinfo: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			return false, fmt.Errorf("userinfo endpoint returned status %d", resp.StatusCode)
		}
		
		// 성공적으로 사용자 정보를 받아왔으면 true 반환
		return true, nil
	}

	// 이후 ID 토큰 검증 (idToken이 실제 ID 토큰인 경우)
	oauth2Service, err := gOauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
	if err != nil {
		return false, fmt.Errorf("unable to create OAuth2 service: %v", err)
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return false, fmt.Errorf("unable to verify id token: %v", err)
	}

	// Check if the token's audience matches your app's client ID for web
	if tokenInfo.Audience == localConfig.GetEnvConfig("googleWebClientID") {
		return true, nil
	}

	return false, fmt.Errorf("id token is not issued by this app")
}

func (GoogleUsecae) GetUserEmail(idToken string) (string, error) {
	ctx := context.Background()

	oauth2Service, err := gOauth2.NewService(ctx, option.WithAPIKey(localConfig.GetEnvConfig("googleApiKey")))
	if err != nil {
		return "", fmt.Errorf("unable to create OAuth2 service: %v", err)
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return "", fmt.Errorf("unable to verify id token: %v", err)
	}

	// Check if the token's audience matches any of your app's client IDs
	if tokenInfo.Audience == localConfig.GetEnvConfig("googleClientID") ||
		tokenInfo.Audience == localConfig.GetEnvConfig("googleAndroidClientID") ||
		tokenInfo.Audience == localConfig.GetEnvConfig("googleWebClientID") {
		
		if tokenInfo.Email != "" {
			return tokenInfo.Email, nil
		}
	}

	return "", fmt.Errorf("id token is not issued by this app or unable to retrieve user's email")
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
