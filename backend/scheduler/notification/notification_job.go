package notification

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"joosum-backend/app/setting"
	"joosum-backend/pkg/config"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// type Notification struct {
// 	NotificationId string    `bson:"_id"`
// 	Title          string    `bson:"title"`
// 	Body           string    `bson:"body"`
// 	LinkCount      int64     `bson:"link_count"`
// 	Type           string    `bson:"type"`
// 	CreatedAt      time.Time `bson:"created_at"`
// 	UserId         string    `bson:"user_id"`
// }

type Message struct {
	Token        string       `json:"token" bson:"token"`
	Notification Notification `json:"notification" bson:"notification"`
}

type Notification struct {
	Body  string `json:"body" bson:"body"`
	Title string `json:"title" bson:"title"`
}

type Res struct {
	Name string `json:"name" bson:"name"`
}

func SendUnreadLink() error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := setting.SettingModel{}

	// 1. device token 가져옴
	s.GetNotificationAgrees()

	// 2. 알림 보냄
	projectId := config.GetEnvConfig("projectId")
	deviceToken := config.GetEnvConfig("deviceToken")

	googleToken, err := getAccesstoken()
	if err != nil {
		return err
	}

	client := resty.New()
	client.SetAuthToken(googleToken)

	msg := Message{
		Token: deviceToken,
		Notification: Notification{
			Title: "title new new 333!!!",
			Body:  "body new new 3333 !!",
		},
	}

	var res Res
	resp, _ := client.R().
		SetBody(map[string]Message{"message": msg}).
		SetResult(&res).
		// SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post("https://fcm.googleapis.com/v1/projects/" + projectId + "/messages:send")

	fmt.Println(resp)

	// 3. 알림 저장

	// setting.SettingModel.GetNotificationAgrees()

	// notification := Notification{
	// 	NotificationId: util.CreateId("Notification"),
	// 	Title:          "읽지 않은 링크가 있어요.",
	// 	Body:           "읽지 않은 링크가 n건 있어요.\n저장해 둔 링크를 확인해보세요!",
	// 	LinkCount:      5,
	// 	Type:           "unread",
	// 	CreatedAt:      time.Now(),
	// 	UserId:         "User-dea95e0a-6d06-4d9f-bd2e-094bcedcc792",
	// }
	// _, err := db.NotificationCollection.InsertOne(ctx, notification)
	// if err != nil {
	// 	return err
	// }

	log.Println("Successfully send notification about unread link!!")

	return nil
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

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"

type tokenProvider struct {
	tokenSource oauth2.TokenSource
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
