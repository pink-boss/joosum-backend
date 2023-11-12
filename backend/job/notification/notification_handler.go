package notification

import (
	"log"
)

const firebaseScope = "https://www.googleapis.com/auth/firebase.messaging"

func SendUnreadLink() {

	// 1. device token 가져옴
	notificationAgrees, err := getNotificationAgrees()
	if err != nil {
		panic("알림동의 목록을 가져오는데 실패했습니다: " + err.Error())
	}
	log.Printf("%d 개의 알림동의 정보를 가져왔습니다.\n\n", len(notificationAgrees))

	// 2. 알림 보냄, 저장
	err = SendUnreadLinks(notificationAgrees)
	if err != nil {
		panic("failed to send or save notifications: " + err.Error())
	}
	log.Println("END")
}

func SendUnclassifiedLink() {

	// 1. device token 가져옴
	notificationAgrees, err := getNotificationAgrees()
	if err != nil {
		panic("알림동의 목록을 가져오는데 실패했습니다: " + err.Error())
	}
	log.Printf("%d 개의 알림동의 정보를 가져왔습니다.\n\n", len(notificationAgrees))

	// 2. 알림 보냄, 저장
	err = SendUnclassifiedLinks(notificationAgrees)
	if err != nil {
		panic("failed to send or save notifications: " + err.Error())
	}
	log.Println("END")
}
