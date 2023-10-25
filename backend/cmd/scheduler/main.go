package main

import (
	"joosum-backend/job/notification"
	"joosum-backend/pkg/config"
	"joosum-backend/pkg/util"
	"log"
	"os"
)

func main() {
	config.EnvSchedulerConfig()
	util.StartMongoDB()

	var notificationType string
	if len(os.Args) > 1 {
		notificationType = os.Args[1]
	}

	if notificationType == "unread" {
		notification.SendUnreadLink()
	} else if notificationType == "unclassified" {
		notification.SendUnclassifiedLink()
	} else {
		log.Fatal("invalid notificationType")
	}

	util.CloseMongoDB()
}
