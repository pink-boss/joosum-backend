package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// userCollection은 User 모델의 컬렉션 인스턴스를 저장합니다.
var UserCollection *mongo.Collection

// InitUserCollection은 전달된 클라이언트 인스턴스를 사용하여 userCollection 변수를 설정합니다.
func InitUserCollection(client *mongo.Client, dbName string) {
	UserCollection = client.Database(dbName).Collection("users")
	EnsureIndexes(UserCollection)
}

// TO DO
// Index 생성, 본인의 Collection 인스턴스 변수, 해당 collection을 init 하는 함수는
// 공통으로 쓰일 것 같으니 패턴화 해서 분리해두는 것이 좋을 것 같습니다.

// email에 대한 인덱스 생성
func EnsureIndexes(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

var LinkBookCollection *mongo.Collection

func InitLinkBookCollection(client *mongo.Client, dbName string) {
	LinkBookCollection = client.Database(dbName).Collection("linkBooks")
}
