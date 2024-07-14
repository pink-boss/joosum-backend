package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

// email에 대한 인덱스 생성
func UserEnsureIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func InitUserCollection(client *mongo.Client, dbName string) {
	UserCollection = client.Database(dbName).Collection("users")
	UserEnsureIndexes(UserCollection)
}

var LinkCollection *mongo.Collection

func LinkEnsureIndexes(collection *mongo.Collection) error {
	userIdIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(false),
	}

	linkIdIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "link_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, userIdIndexModel)

	if err != nil {
		return err
	}

	_, err = collection.Indexes().CreateOne(ctx, linkIdIndexModel)

	return err
}

func InitLinkCollection(client *mongo.Client, dbName string) {
	LinkCollection = client.Database(dbName).Collection("links")
	LinkEnsureIndexes(UserCollection)
}

var LinkBookCollection *mongo.Collection

func InitLinkBookCollection(client *mongo.Client, dbName string) {
	LinkBookCollection = client.Database(dbName).Collection("linkBooks")
}

var InactiveUserCollection *mongo.Collection

func InactiveUserEnsureIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func InitInactiveUserCollection(client *mongo.Client, dbName string) {
	InactiveUserCollection = client.Database(dbName).Collection("inactiveUsers")
	InactiveUserEnsureIndexes(InactiveUserCollection)
}

var TagCollection *mongo.Collection

// InitUserCollection은 전달된 클라이언트 인스턴스를 사용하여 userCollection 변수를 설정합니다.
func InitTagCollection(client *mongo.Client, dbName string) {
	TagCollection = client.Database(dbName).Collection("tags")
	EnsureIndexes(TagCollection)
}

// TO DO
// Index 생성, 본인의 Collection 인스턴스 변수, 해당 collection을 init 하는 함수는
// 공통으로 쓰일 것 같으니 패턴화 해서 분리해두는 것이 좋을 것 같습니다.

// email에 대한 인덱스 생성
func EnsureIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "tag_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

var NotificationCollection *mongo.Collection

func InitNotificationCollection(client *mongo.Client, dbName string) {
	NotificationCollection = client.Database(dbName).Collection("notifications")
}

var NotificationAgreeCollection *mongo.Collection

func InitNotificationAgreeCollection(client *mongo.Client, dbName string) {
	NotificationAgreeCollection = client.Database(dbName).Collection("notificationAgrees")
}

var BannerCollection *mongo.Collection

func InitBannerCollection(client *mongo.Client, dbName string) {
	BannerCollection = client.Database(dbName).Collection("banners")
}
