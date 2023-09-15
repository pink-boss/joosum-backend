package notification

import (
	"context"
	"joosum-backend/pkg/db"
	"joosum-backend/pkg/util"
	"log"
	"time"
)

type Notification struct {
	NotificationId string    `bson:"_id"`
	Title          string    `bson:"title"`
	Body           string    `bson:"body"`
	LinkCount      int64     `bson:"link_count"`
	Type           string    `bson:"type"`
	CreatedAt      time.Time `bson:"created_at"`
	UserId         string    `bson:"user_id"`
}

func SendUnreadLink() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notification := Notification{
		NotificationId: util.CreateId("Notification"),
		Title:          "읽지 않은 링크가 있어요.",
		Body:           "읽지 않은 링크가 n건 있어요.\n저장해 둔 링크를 확인해보세요!",
		LinkCount:      5,
		Type:           "unread",
		CreatedAt:      time.Now(),
		UserId:         "User-dea95e0a-6d06-4d9f-bd2e-094bcedcc792",
	}
	_, err := db.NotificationCollection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}

	log.Println("Successfully send notification about unread link!!")

	return nil
}

func SendUnclassifiedLink() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notification := Notification{
		NotificationId: util.CreateId("Notification"),
		Title:          "분류되지 않은 링크가 있어요.",
		Body:           "분류되지 않은 링크가 n건 있어요. \n폴더를 만들어서 정리해보세요!",
		LinkCount:      5,
		Type:           "unclassified",
		CreatedAt:      time.Now(),
		UserId:         "User-dea95e0a-6d06-4d9f-bd2e-094bcedcc792",
	}
	_, err := db.NotificationCollection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}

	log.Println("Successfully send notification about unclassified link!!")

	return nil
}
