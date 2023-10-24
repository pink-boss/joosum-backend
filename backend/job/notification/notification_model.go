package notification

import (
	"time"

	"golang.org/x/oauth2"
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

type tokenProvider struct {
	tokenSource oauth2.TokenSource
}
