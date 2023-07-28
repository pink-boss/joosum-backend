package util

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
	db.InitCommonCollection(client, dbName)
}

func LoadCommonData() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := db.CommonCollection.FindOne(ctx, bson.M{"type": "DEFAULT_FOLDER"}).Decode(&db.DefaultFolder)
	if err != nil {
		log.Fatalf("Failed to get the DEFAULT_FOLDER : %v", err)
	}
}
