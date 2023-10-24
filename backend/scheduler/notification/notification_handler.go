package notification

import (
	"joosum-backend/app/setting"
	"log"
)

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"

func SendUnreadLink() {
	setting := setting.SettingModel{}

	// 1. device token 가져옴
	notificationAgrees, err := setting.GetNotificationAgrees("unread")
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
	notificationAgrees, err := setting.GetNotificationAgrees("unclassified")
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
