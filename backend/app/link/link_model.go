package link

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Link struct {
	LinkId     string    `bson:"link_id" json:"linkId"`
	URL        string    `bson:"url" json:"url"`
	UserID     string    `bson:"user_id" json:"userId"`
	Title      string    `bson:"title" json:"title"`
	LinkBookId string    `bson:"link_book_id" json:"linkBookId"`
	ReadCount  int       `bson:"read_count" json:"readCount"`
	LastReadAt time.Time `bson:"last_read_at" json:"LastReadAt"`
	CreatedAt  time.Time `bson:"created_at" json:"CreatedAt"`
	UpdatedAt  time.Time `bson:"updated_at" json:"UpdatedAt"`
}

type LinkModel struct {
}

var linkCollection *mongo.Collection

func InitLinkCollection(client *mongo.Client, dbName string) {
	linkCollection = client.Database(dbName).Collection("links")
	EnsureIndexes(linkCollection)
}

func EnsureIndexes(collection *mongo.Collection) error {
	userIdIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
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

func (LinkModel) CreateLink(url string, title string, userId string, linkBookId string) (*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uniqueKey := make(chan string)

	// TO DO : uid generater 만들기
	go func() {
		s := "Link-"
		buf := bytes.NewBufferString(s)
		uid := uuid.New()
		buf.WriteString(uid.String())
		uniqueKey <- buf.String()
	}()

	uid := <-uniqueKey

	link := Link{
		LinkId:     uid,
		URL:        url,
		Title:      title,
		UserID:     userId,
		LinkBookId: linkBookId,
		ReadCount:  0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := linkCollection.InsertOne(ctx, link)

	return &link, err
}

func (LinkModel) GetAllLinkByUserId(userId string) ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	cursor, err := linkCollection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var link Link
		cursor.Decode(&link)
		links = append(links, &link)
	}

	return links, nil
}

func (LinkModel) GetAllLinkByUserIdAndLinkBookId(userId string, linkBookId string) ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	cursor, err := linkCollection.Find(ctx, bson.M{"user_id": userId, "link_book_id": linkBookId})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var link Link
		cursor.Decode(&link)
		links = append(links, &link)
	}

	return links, nil
}

func (LinkModel) GetOneLinkByLinkId(linkId string) (*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var link *Link

	error := linkCollection.FindOne(ctx, bson.M{"link_id": linkId}).Decode(&link)

	if error != nil {
		return nil, error
	}

	return link, nil
}

func (LinkModel) DeleteOneByLinkId(linkId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := linkCollection.DeleteOne(ctx, bson.M{"link_id": linkId})

	return err
}

func (LinkModel) DeleteAllLinksByUserId(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := linkCollection.DeleteMany(ctx, bson.M{"user_id": userId})

	return err
}

func (LinkModel) DeleteAllLinksByLinkBookId(userId string, linkBookId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := linkCollection.DeleteMany(ctx, bson.M{"user_id": userId, "link_book_id": linkBookId})

	return err
}

func (LinkModel) UpdateReadCountByLinkId(linkId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := linkCollection.UpdateOne(ctx, bson.M{"link_id": linkId}, bson.M{"last_read_at": time.Now(), "$inc": bson.M{"read_count": 1}})

	return err
}

func (LinkModel) UpdateLinkBookIdByLinkId(linkId string, linkBookId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := linkCollection.UpdateOne(ctx, bson.M{"link_id": linkId}, bson.M{"link_book_id": linkBookId})

	return err
}

func (LinkModel) UpdateTitleAndUrlByLinkId(linkId string, url string, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := linkCollection.UpdateOne(ctx, bson.M{"link_id": linkId}, bson.M{"url": url, "title": title})

	return err
}
