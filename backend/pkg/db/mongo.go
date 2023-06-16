package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// clientInstance는 MongoDB 클라이언트의 싱글톤 인스턴스입니다.
var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

// GetMongoClient는 주어진 MongoDB URI를 사용하여 싱글톤 MongoDB 클라이언트 인스턴스를 반환합니다.
// 이미 인스턴스가 있다면 기존 인스턴스를 반환하고, 그렇지 않으면 새 인스턴스를 생성한 후 반환합니다.
func GetMongoClient(ctx context.Context, mongoURI string) (*mongo.Client, error) {
	var err error

	clientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(mongoURI)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
			return
		}
		if client == nil {
			log.Fatal("Mongo client is nil")
		}
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		clientInstance = client
	})

	return clientInstance, err
}

// DisconnectMongoClient는 싱글톤 MongoDB 클라이언트 인스턴스의 연결을 종료합니다.
func DisconnectMongoClient(ctx context.Context) error {
	if clientInstance != nil {
		return clientInstance.Disconnect(ctx)
	}
	return nil
}
