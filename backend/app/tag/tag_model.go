package tag

import (
	"bytes"
	"context"
	"joosum-backend/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type TagModel struct{}

type Tag struct {
	ID     string   `json:"-"`
	Names  []string `json:"names"`
	UserId string   `json:"user_id" bson:"user_id"`
}

func (TagModel) UpsertTags(userId string, names []string) (*Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uniqueKey := make(chan string)

	// TO DO : uid generater 만들기
	go func() {
		s := "Tag-"
		buf := bytes.NewBufferString(s)
		uid := uuid.New()
		buf.WriteString(uid.String())
		uniqueKey <- buf.String()
	}()

	uid := <-uniqueKey

	tags := Tag{
		ID:     uid,
		Names:  names,
		UserId: userId,
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"user_id": userId}
	update := bson.M{
		"$setOnInsert": bson.M{"_id": uid}, // update 할 때 id 가 바뀌지 않도록 함
		"$set":         bson.M{"names": names},
	}

	_, err := db.TagCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return &tags, nil
}

func (TagModel) FindTagsByUserId(userId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}
	var tags = &Tag{}

	err := db.TagCollection.FindOne(ctx, filter).Decode(tags)
	if err != nil {
		return nil, err
	}

	return tags.Names, nil
}

func (TagModel) DeleteTag(user_id string, tag_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": user_id,
		"tag_id":  tag_id,
	}

	_, err := db.TagCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
