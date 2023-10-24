package setting

import (
	"context"
	"joosum-backend/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SettingModel struct {
}

type NotificationAgree struct {
	NotificationAgreeId *string `bson:"_id" json:"notificationAgreeId" example:"652bf9508de1187ff1e16e24"`
	DeviceId            *string `bson:"device_id" json:"deviceId"`
	IsReadAgree         bool    `bson:"is_read_agree" json:"isReadAgree"`
	IsClassifyAgree     bool    `bson:"is_classify_agree" json:"isClassifyAgree"`
	UserId              string  `bson:"user_id" json:"userId" example:"User-dea95e0a-6d06-4d9f-bd2e-094bcedcc792"`
}

type DeviceReq struct {
	DeviceId string `json:"deviceId"`
}

type PushNotificationReq struct {
	IsReadAgree     bool `json:"isReadAgree"`
	IsClassifyAgree bool `json:"isClassifyAgree"`
}

func (SettingModel) SaveDeviceId(deviceId, userId string) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"user_id", userId}}
	update := bson.D{{"$set", bson.D{{"device_id", deviceId}}}}
	opts := options.Update().SetUpsert(true)

	result, err := db.NotificationAgreeCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (SettingModel) GetNotificationAgree(userId string) (*NotificationAgree, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"user_id", userId}}
	var agree NotificationAgree
	err := db.NotificationAgreeCollection.FindOne(ctx, filter).Decode(&agree)
	if err != nil {
		return nil, err
	}

	return &agree, nil
}

func (SettingModel) UpdatePushNotification(req PushNotificationReq, userId string) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"user_id", userId}}
	update := bson.D{{"$set", bson.D{
		{"is_read_agree", req.IsReadAgree},
		{"is_classify_agree", req.IsClassifyAgree},
	}}}
	opts := options.Update().SetUpsert(true)

	result, err := db.NotificationAgreeCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (SettingModel) DeleteDivceId(userId string) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"user_id", userId}}
	update := bson.D{{"$set", bson.D{
		{"device_id", nil},
	}}}
	opts := options.Update().SetUpsert(true)

	result, err := db.NotificationAgreeCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil

}
