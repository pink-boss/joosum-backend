package notif

import (
	"context"
	"joosum-backend/pkg/db"
	"time"

	. "github.com/gobeam/mongo-go-pagination"
	mongopagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationModel struct {
}

type Notification struct {
	NotificationId string    `json:"notificationId" bson:"_id" example:"Notification-bb61159e-da11-47b2-bdfc-fbda0e85a856"`
	Title          string    `json:"title" bson:"title" example:"읽지 않은 링크가 3건 있어요."`
	Body           string    `json:"body" bson:"body" example:"저장해 둔 링크를 확인해보세요!"`
	IsRead         bool      `json:"isRead" bson:"is_read"`
	Type           string    `json:"type" bson:"type" example:"unread"`
	CreatedAt      time.Time `json:"createdAt" bson:"created_at"`
	UserId         string    `json:"userId" bson:"user_id" example:"User-590e39b3-7661-4387-8501-85aaf87d133c"`
}

type NotificationRes struct {
	Notifications []Notification                  `json:"notifications" bson:"notifications"`
	Page          *mongopagination.PaginationData `json:"page" bson:"page"`
}

type notificationResDocs struct {
	Notifications []Notification     `json:"notifications" bson:"notifications"`
	Page          *db.PaginationData `json:"page" bson:"page"`
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

func (NotificationModel) UpdateIsRead(notificationId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_read": true,
		},
	}

	err := db.NotificationCollection.FindOneAndUpdate(ctx, bson.M{"_id": notificationId}, update).Decode(&mongo.SingleResult{})
	if err != nil {
		return err
	}
	return nil
}
