package notification

import (
	"context"
	"joosum-backend/app/setting"
	"joosum-backend/pkg/db"
	"joosum-backend/pkg/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
)

// 동의 타입
const (
	Unread       = "unread"
	Unclassified = "unclassified"
)

type Notification struct {
	NotificationId string    `bson:"_id"`
	Title          string    `bson:"title"`
	Body           string    `bson:"body"`
	IsRead         bool      `bson:"is_read"`
	Type           string    `bson:"type"`
	CreatedAt      time.Time `bson:"created_at"`
	UserId         string    `bson:"user_id"`
}

type FcmReq struct {
	Token        string          `json:"token" bson:"token"`
	Notification FcmNotification `json:"notification" bson:"notification"`
}

type FcmNotification struct {
	Body  string `json:"body" bson:"body"`
	Title string `json:"title" bson:"title"`
}

type FcmRes struct {
	Name string `json:"name" bson:"name"`
}

type NotificationResult struct {
	UserId string
	Msg    string
	Err    error
}

type tokenProvider struct {
	tokenSource oauth2.TokenSource
}

func SaveNotification(userId, title, body, notificationType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notification := Notification{
		NotificationId: util.CreateId("Notification"),
		Title:          title,
		Body:           body,
		Type:           notificationType,
		CreatedAt:      time.Now(),
		UserId:         userId,
	}

	_, err := db.NotificationCollection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}
	return nil
}

func getNotificationAgrees() ([]setting.NotificationAgree, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := db.NotificationAgreeCollection.Find(ctx, bson.M{})

	var results []setting.NotificationAgree
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
