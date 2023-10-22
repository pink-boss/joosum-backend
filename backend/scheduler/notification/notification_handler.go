package notification

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"joosum-backend/app/setting"
	"log"
	"time"

	"golang.org/x/oauth2/google"
)

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"

func SendUnreadLink() {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	setting := setting.SettingModel{}

	// 1. device token 가져옴
	notificationAgrees, err := setting.GetNotificationAgrees()
	if err != nil {
		panic("failed to get the device tokens: " + err.Error())
	}

	// 2. 알림 보냄, 저장
	err = sendAndSaveNotifications(notificationAgrees)
	if err != nil {
		panic("failed to send or save notifications: " + err.Error())
	}

	log.Println("Successfully send notification about unread link!!")
}

func SendUnclassifiedLink() error {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// notification := Notification{
	// 	NotificationId: util.CreateId("Notification"),
	// 	Title:          "분류되지 않은 링크가 있어요.",
	// 	Body:           "분류되지 않은 링크가 n건 있어요. \n폴더를 만들어서 정리해보세요!",
	// 	LinkCount:      5,
	// 	Type:           "unclassified",
	// 	CreatedAt:      time.Now(),
	// 	UserId:         "User-dea95e0a-6d06-4d9f-bd2e-094bcedcc792",
	// }
	// _, err := db.NotificationCollection.InsertOne(ctx, notification)
	// if err != nil {
	// 	return err
	// }

	log.Println("Successfully send notification about unclassified link!!")

	return nil
}

func getAccesstoken() (string, error) {
	tokenProvider, err := newTokenProvider("fireBaseKey.json")
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
		return nil, errors.New("fcm: failed to get JWT config for the firebase.messaging scope")
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
