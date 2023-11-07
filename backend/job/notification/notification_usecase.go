package notification

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"joosum-backend/app/link"
	"joosum-backend/app/setting"
	"joosum-backend/pkg/config"
	"log"
	. "strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/google"
)

// count == 0 알림 x, 저장 x
// device id == null 알림 x 저장 o
// 동의여부 == true 면 알림 o 저장 o, false 면 알림 x 저장 o
func SendUnreadLinks(notificationAgrees []setting.NotificationAgree) error {
	linkModel := link.LinkModel{}

	projectId := config.GetEnvConfig("projectId")
	googleToken, err := getAccesstoken()
	if err != nil {
		return err
	}
	var successUserIds []string
	var failUserIds []NotificationResult

	client := resty.New()
	client.SetAuthToken(googleToken)
	for _, notificationAgree := range notificationAgrees {
		title := "읽지 않은 링크가 n건 있어요."
		body := "저장해 둔 링크를 확인해보세요!"

		userId := notificationAgree.UserId
		deviceToken := notificationAgree.DeviceId

		// 읽지않은 링크 갯수 세기
		unreadLinkCnt, err := linkModel.GetUserUnreadLinkCount(userId)
		if err != nil {
			failUserIds = append(failUserIds, NotificationResult{userId, "링크갯수 조회 실패", err})
			continue
		}

		// count == 0 이면 패스
		if unreadLinkCnt == 0 {
			continue
		}

		title = strings.Replace(title, "n", FormatInt(unreadLinkCnt, 10), 1)

		// device id 가 null 이 아니고, 알림 동의 일 때
		if deviceToken != nil && notificationAgree.IsReadAgree {

			// Firebase 로 알림 보내기
			result := sendNotification(client, deviceToken, title, body, projectId)

			if result.IsSuccess() {
				successUserIds = append(successUserIds, userId)

			} else {
				failUserIds = append(failUserIds, NotificationResult{userId, "알림발송 실패", errors.New("")})
			}
		}

		// 알림 저장
		err = SaveNotification(userId, title, body, Unread)
		if err != nil {
			failUserIds = append(failUserIds, NotificationResult{userId, "알림저장 실패", err})
		}
	}

	log.Println("\t\t[알림발송 실패 목록]\n")
	log.Printf("               UserId                           Message                       Error")
	for _, failUser := range failUserIds {
		log.Printf("%s \t %s \t %s", failUser.UserId, failUser.Msg, failUser.Err)
	}
	log.Println()
	log.Printf("successUserIds=%v \n\n %d 개의 '읽지않은 링크' 알림을 보내는데 성공했습니다.\n\n", successUserIds, len(successUserIds))

	return nil
}

func SendUnclassifiedLinks(notificationAgrees []setting.NotificationAgree) error {
	linkModel := link.LinkModel{}
	linkBookModel := link.LinkBookModel{}

	projectId := config.GetEnvConfig("projectId")
	googleToken, err := getAccesstoken()
	if err != nil {
		return err
	}
	var successUserIds []string
	var failUserIds []NotificationResult

	client := resty.New()
	client.SetAuthToken(googleToken)
	for _, notificationAgree := range notificationAgrees {
		title := "분류되지 않은 링크가 n건 있어요."
		body := "폴더를 만들어서 정리해보세요!"

		userId := notificationAgree.UserId
		deviceToken := notificationAgree.DeviceId

		// 분류되지 않은 링크 갯수 세기
		defaultLinkBook, err := linkBookModel.GetDefaultLinkBook(userId)
		if err != nil {
			failUserIds = append(failUserIds, NotificationResult{userId, "기본폴더 조회 실패", err})
			continue
		}

		unclassifyCnt, err := linkModel.GetLinkBookLinkCount(defaultLinkBook.LinkBookId)
		if err != nil {
			failUserIds = append(failUserIds, NotificationResult{userId, "링크갯수 조회 실패", err})
			continue
		}

		// count == 0 이면 패스
		if unclassifyCnt == 0 {
			continue
		}

		title = strings.Replace(title, "n", FormatInt(unclassifyCnt, 10), 1)

		// device id 가 null 이 아니고, 알림 동의 일 때
		if deviceToken != nil && notificationAgree.IsClassifyAgree {

			// Firebase 로 알림 보내기
			result := sendNotification(client, deviceToken, title, body, projectId)

			if result.IsSuccess() {
				successUserIds = append(successUserIds, userId)

			} else {
				failUserIds = append(failUserIds, NotificationResult{userId, "알림발송 실패", errors.New("")})
			}
		}

		// 알림 저장
		err = SaveNotification(userId, title, body, Unclassified)
		if err != nil {
			failUserIds = append(failUserIds, NotificationResult{userId, "알림저장 실패", err})
		}
	}

	log.Println("\t\t[알림발송 실패 목록]\n")
	log.Printf("               UserId                           Message                       Error")
	for _, failUser := range failUserIds {
		log.Printf("%s \t %s \t %s", failUser.UserId, failUser.Msg, failUser.Err)
	}
	log.Println()
	log.Printf("successUserIds=%v \n\n %d 개의 '분류되지 않은 링크' 알림을 보내는데 성공했습니다.\n\n", successUserIds, len(successUserIds))

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

func sendNotification(client *resty.Client, deviceToken *string, title, body, projectId string) *resty.Response {
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

	return result
}
