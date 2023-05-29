package tag

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TagModel struct {}

type Tag struct{
	ID string `json:"tag_id"`
	Name string `json:"name"`
	UserId string `json:"user_id"`
}

var tagCollection *mongo.Collection

// InitUserCollection은 전달된 클라이언트 인스턴스를 사용하여 userCollection 변수를 설정합니다.
func InitUserCollection(client *mongo.Client, dbName string) {
	tagCollection = client.Database(dbName).Collection("tags")
	EnsureIndexes(tagCollection)
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

func (TagModel)CreateTag(name string, user_id string) (*Tag, error){
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

	tag := &Tag{
		ID: uid,
		Name: name,
		UserId: user_id,
	}

	_, err := tagCollection.InsertOne(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (TagModel)FindTagByUserId(user_id string) ([]*Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": user_id}
	tags := []*Tag{}

	cursor, err := tagCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func (TagModel)DeleteTag(user_id string, tag_id string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": user_id,
		"tag_id": tag_id,
	}

	_, err := tagCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}