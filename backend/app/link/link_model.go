package link

import (
	"bytes"
	"context"
	"joosum-backend/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Link struct {
	LinkId       string    `bson:"link_id" json:"linkId"`
	URL          string    `bson:"url" json:"url"`
	UserID       string    `bson:"user_id" json:"userId"`
	Title        string    `bson:"title" json:"title"`
	LinkBookId   string    `bson:"link_book_id" json:"linkBookId"`
	LinkBookName string    `bson:"link_book_name" json:"linkBookName"`
	ThumbnailURL string    `bson:"thumbnail_url" json:"thumbnailURL"`
	Tags         []string  `bson:"tags" json:"tags"`
	ReadCount    int       `bson:"read_count" json:"readCount"`
	LastReadAt   time.Time `bson:"last_read_at" json:"lastReadAt"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
}

type LinkModel struct {
}

func (LinkModel) CreateLink(url string, title string, userId string, linkBookId string, linkBookName string, thumbnailURL string, tags []string) (*Link, error) {
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
		LinkId:       uid,
		URL:          url,
		Title:        title,
		UserID:       userId,
		LinkBookId:   linkBookId,
		LinkBookName: linkBookName,
		ThumbnailURL: thumbnailURL,
		Tags:         tags,
		ReadCount:    0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := db.LinkCollection.InsertOne(ctx, link)

	return &link, err
}

func (LinkModel) Get9LinksByUserId(userId string) ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	cursor, err := db.LinkCollection.Find(ctx, bson.M{"user_id": userId}, options.Find().SetLimit(9))
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

func (LinkModel) GetAllLinkByUserId(userId string) ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	cursor, err := db.LinkCollection.Find(ctx, bson.M{"user_id": userId})
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

	cursor, err := db.LinkCollection.Find(ctx, bson.M{"user_id": userId, "link_book_id": linkBookId})
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

	error := db.LinkCollection.FindOne(ctx, bson.M{"link_id": linkId}).Decode(&link)

	if error != nil {
		return nil, error
	}

	return link, nil
}

func (LinkModel) GetLinkBookLinkCount(linkBookId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, error := db.LinkCollection.CountDocuments(ctx, bson.M{"link_book_id": linkBookId})
	if error != nil {
		return 0, error
	}

	return result, nil
}

func (LinkModel) GetUserLinkCount(userId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, error := db.LinkCollection.CountDocuments(ctx, bson.M{"user_id": userId})
	if error != nil {
		return 0, error
	}

	return result, nil
}

func (LinkModel) DeleteOneByLinkId(linkId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.LinkCollection.DeleteOne(ctx, bson.M{"link_id": linkId})

	return err
}

func (LinkModel) DeleteAllLinksByUserId(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.LinkCollection.DeleteMany(ctx, bson.M{"user_id": userId})

	return err
}

func (LinkModel) DeleteAllLinksByLinkBookId(userId string, linkBookId string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.LinkCollection.DeleteMany(ctx, bson.M{"user_id": userId, "link_book_id": linkBookId})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (LinkModel) UpdateReadCountByLinkId(linkId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.LinkCollection.UpdateOne(ctx, bson.M{"link_id": linkId}, bson.M{"$set": bson.M{"last_read_at": time.Now()}, "$inc": bson.M{"read_count": 1}})

	return err
}

func (LinkModel) UpdateLinkBookIdAndTitleByLinkId(linkId string, linkBookId string, linkBookName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.LinkCollection.UpdateOne(ctx, bson.M{"link_id": linkId}, bson.M{"$set": bson.M{"link_book_id": linkBookId, "link_book_name": linkBookName}})

	return err
}

func (LinkModel) UpdateTitleAndUrlByLinkId(linkId string, url string, title string, thumbnailURL string, tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateFields := bson.M{}

	if url != "" {
		updateFields["url"] = url
	}

	if title != "" {
		updateFields["title"] = title
	}

	if thumbnailURL != "" {
		updateFields["thumbnailURL"] = thumbnailURL
	}

	if len(tags) != 0 {
		updateFields["tags"] = tags
	}

	if len(updateFields) == 0 {
		return nil
	}

	_, err := db.LinkCollection.UpdateOne(ctx, bson.M{"link_id": linkId}, bson.M{"$set": updateFields})

	return err
}
