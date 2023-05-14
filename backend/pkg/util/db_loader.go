package util

import (
	"context"
	"joosum-backend/app/user"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/database"
	"log"
	"time"
)

func StartMongoDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := database.GetMongoClient(ctx, config.GetEnvConfig("mongoDB"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Collection load
	user.InitUserCollection(client, config.GetEnvConfig("dbName"))

}
