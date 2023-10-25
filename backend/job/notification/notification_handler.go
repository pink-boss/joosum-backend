package notification

import (
	"log"
)

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"

func SendUnreadLink() {

	// 1. device token 가져옴
	notificationAgrees, err := getNotificationAgrees()
	if err != nil {
		panic("failed to get the device tokens: " + err.Error())
	}

	// 2. 알림 보냄, 저장
	err = SendUnreadLinks(notificationAgrees)
	if err != nil {
		panic("failed to send or save notifications: " + err.Error())
	}

	log.Println("Successfully send notification about unread link!!")
}

func SendUnclassifiedLink() {

	// 1. device token 가져옴
	notificationAgrees, err := getNotificationAgrees()
	if err != nil {
		panic("failed to get the device tokens: " + err.Error())
	}

	// 2. 알림 보냄, 저장
	err = SendUnclassifiedLinks(notificationAgrees)
	if err != nil {
		panic("failed to send or save notifications: " + err.Error())
	}

	log.Println("Successfully send notification about unclassified link!!")
}
