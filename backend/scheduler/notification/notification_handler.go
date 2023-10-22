package notification

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"joosum-backend/app/setting"
	"log"

	"golang.org/x/oauth2/google"
)

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"

func SendUnreadLink() {
	setting := setting.SettingModel{}

	// 1. device token 가져옴
	notificationAgrees, err := setting.GetNotificationAgrees()
	if err != nil {
		panic("failed to get the device tokens: " + err.Error())
	}

	// 2. 알림 보냄, 저장
	err = sendAndSaveNotifications(notificationAgrees, "unread")
	if err != nil {
		panic("failed to send or save notifications: " + err.Error())
	}

	log.Println("Successfully send notification about unread link!!")
}

func SendUnclassifiedLink() {
	setting := setting.SettingModel{}

	// 1. device token 가져옴
	notificationAgrees, err := setting.GetNotificationAgrees()
	if err != nil {
		panic("failed to get the device tokens: " + err.Error())
	}

	// 2. 알림 보냄, 저장
	err = sendAndSaveNotifications(notificationAgrees, "unclassified")
	if err != nil {
		panic("failed to send or save notifications: " + err.Error())
	}

	log.Println("Successfully send notification about unclassified link!!")
}

func getAccesstoken() (string, error) {
	tokenProvider, err := newTokenProvider("../../fireBaseKey.json")
	if err != nil {
		return "", fmt.Errorf("Failed to get Token provider: %v", err)
	}
	token, err := tokenProvider.token()
	if err != nil {
		return "", fmt.Errorf("Failed to get Token: %v", err)
	}

	return token, nil
}

// newTokenProvider function to get token for fcm-send
func newTokenProvider(credentialsLocation string) (*tokenProvider, error) {
	jsonKey, err := ioutil.ReadFile(credentialsLocation)
	if err != nil {
		return nil, errors.New("fcm: failed to read credentials file at: " + credentialsLocation)
	}
	cfg, err := google.JWTConfigFromJSON(jsonKey, firebaseScope)
	if err != nil {
		return nil, errors.New("fcm: failed to get JWT config for the firebase.messaging scope: " + err.Error())
	}
	ts := cfg.TokenSource(context.Background())
	return &tokenProvider{
		tokenSource: ts,
	}, nil
}

func (src *tokenProvider) token() (string, error) {
	token, err := src.tokenSource.Token()
	if err != nil {
		return "", errors.New("fcm: failed to generate Bearer token")
	}
	return token.AccessToken, nil
}
