package notification

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"joosum-backend/app/link"
	"joosum-backend/app/setting"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/db"
	"joosum-backend/pkg/util"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/google"
)

func sendAndSaveNotifications(notificationAgrees []setting.NotificationAgree, notificationType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	projectId := config.GetEnvConfig("projectId")
	googleToken, err := getAccesstoken()
	if err != nil {
		return err
	}
	var successUserIds []string
	var failUserIds []string

	client := resty.New()
	client.SetAuthToken(googleToken)
	for _, notificationAgree := range notificationAgrees {
		deviceToken := notificationAgree.DeviceId
		userId := notificationAgree.UserId

		title, body, err := getNotificationText(notificationType, userId)
		if err != nil {
			return err
		}

		// 2. 알림 보내기
		msg := FcmReq{
			Token: *deviceToken,
			Notification: FcmNotification{
				Title: title,
				Body:  body,
			},
		}

		var res FcmRes
		var authErr resty.Request
		result, _ := client.R().
			SetBody(map[string]FcmReq{"message": msg}).
			SetResult(&res).
			SetError(&authErr). // or SetError(AuthError{}).
			Post("https://fcm.googleapis.com/v1/projects/" + projectId + "/messages:send")

		notification := Notification{
			NotificationId: util.CreateId("Notification"),
			Title:          title,
			Body:           body,
			Type:           notificationType,
			CreatedAt:      time.Now(),
			UserId:         userId,
		}

		if result.IsSuccess() {
			successUserIds = append(successUserIds, userId)

			_, err = db.NotificationCollection.InsertOne(ctx, notification)
			if err != nil {
				return err
			}
		} else {
			failUserIds = append(failUserIds, userId)
		}
	}
	log.Printf("successUserIds=%v \n", successUserIds)
	log.Printf("failUserIds=%v \n", failUserIds)

	return nil
}

func getNotificationText(notificationType, userId string) (title, body string, err error) {
	linkBookModel := link.LinkBookModel{}
	linkModel := link.LinkModel{}

	// 읽지않은 링크 갯수 세기
	if notificationType == "unread" {
		unreadLinkCnt, err := linkModel.GetUserUnreadLinkCount(userId)
		if err != nil {
			return "", "", err
		}
		title = fmt.Sprintf("읽지 않은 링크가 %d건 있어요.", unreadLinkCnt)
		body = "저장해 둔 링크를 확인해보세요!"

		// 분류되지 않은 링크 갯수 세기
	} else {
		defaultLinkBook, err := linkBookModel.GetDefaultLinkBook(userId)
		if err != nil {
			return "", "", err
		}
		unclassifyCnt, err := linkModel.GetLinkBookLinkCount(defaultLinkBook.LinkBookId)
		if err != nil {
			return "", "", err
		}

		title = fmt.Sprintf("분류되지 않은 링크가 %d건 있어요.", unclassifyCnt)
		body = "폴더를 만들어서 정리해보세요!"
	}
	return
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
