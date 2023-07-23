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
