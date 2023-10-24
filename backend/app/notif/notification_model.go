package notif

import (
	"context"
	"joosum-backend/pkg/db"
	"time"

	. "github.com/gobeam/mongo-go-pagination"
	mongopagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson"
)

type NotificationModel struct {
}

type Notification struct {
	NotificationId string    `json:"notificationId" bson:"notification_id"`
	Title          string    `json:"title" bson:"title"`
	Body           string    `json:"body" bson:"body"`
	IsRead         bool      `json:"isRead" bson:"is_read"`
	Type           string    `json:"type" bson:"type"`
	CreatedAt      time.Time `json:"createdAt" bson:"created_at"`
	UserId         string    `json:"userId" bson:"user_id"`
}

type NotificationRes struct {
	Notifications []Notification                  `json:"notifications" bson:"notifications"`
	Page          *mongopagination.PaginationData `json:"page" bson:"page"`
}

func (NotificationModel) Notifications(userId string, page int64) (*NotificationRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	filter := bson.M{
		"user_id":    userId,
		"created_at": bson.M{"$gte": thirtyDaysAgo},
	}

	var notifications []Notification
	paginatedData, err := New(db.NotificationCollection).Context(ctx).Limit(10).Page(page).Sort("created_at", -1).Select(bson.D{}).Filter(filter).Decode(&notifications).Find()
	if err != nil {
		panic(err)
	}

	return &NotificationRes{Notifications: notifications, Page: &paginatedData.Pagination}, nil
}
