package util

import (
	"context"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/db"
	"log"
	"time"
)

func StartMongoDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := db.GetMongoClient(ctx, config.GetEnvConfig("mongoDB"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	dbName := config.GetEnvConfig("dbName")

	// Collection load
	db.InitUserCollection(client, dbName)
	db.InitLinkCollection(client, dbName)
	db.InitLinkBookCollection(client, dbName)
	db.InitInactiveUserCollection(client, dbName)
	db.InitTagCollection(client, dbName)

}

func CloseMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.DisconnectMongoClient(ctx)
	if err != nil {
		log.Fatalf("Failed to disconnect to MongoDB: %v", err)
	}
}
