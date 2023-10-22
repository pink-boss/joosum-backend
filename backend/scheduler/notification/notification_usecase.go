package notification

import (
	"context"
	"fmt"
	"joosum-backend/app/link"
	"joosum-backend/app/setting"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/db"
	"joosum-backend/pkg/util"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

func sendAndSaveNotifications(notificationAgrees []setting.NotificationAgree) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	link := link.LinkModel{}

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

		// 1. 읽지않은 링크 갯수 세기
		unreadLinkCnt, err := link.GetUserUnreadLinkCount(userId)
		if err != nil {
			return err
		}

		// 2. 알림 보내기
		title := fmt.Sprintf("읽지 않은 링크가 %d건 있어요.", unreadLinkCnt)
		body := "저장해 둔 링크를 확인해보세요!"

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

		if result.IsSuccess() {
			successUserIds = append(successUserIds, userId)
		} else {
			failUserIds = append(failUserIds, userId)
		}

		notification := Notification{
			NotificationId: util.CreateId("Notification"),
			Title:          title,
			Body:           body,
			Type:           "unread",
			CreatedAt:      time.Now(),
			UserId:         userId,
		}

		_, err = db.NotificationCollection.InsertOne(ctx, notification)
		if err != nil {
			return err
		}
	}
	log.Printf("successUserIds=%v \n", successUserIds)
	log.Printf("failUserIds=%v \n", failUserIds)

	return nil
}
