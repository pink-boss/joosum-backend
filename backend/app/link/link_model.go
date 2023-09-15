package link

import (
	"bytes"
	"context"
	"joosum-backend/pkg/db"
	"regexp"
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

func (LinkModel) GetAllLink() ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	cursor, err := db.LinkCollection.Find(ctx, bson.M{})
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

func (LinkModel) GetAllLinkByUserId(userId string, sort string) ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	opts := options.Find()
	opts.SetSort(bson.D{
		{Key: sort, Value: 1},
	})

	cursor, err := db.LinkCollection.Find(ctx, bson.M{"user_id": userId}, opts)
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

func (LinkModel) GetAllLinkByUserIdAndSearch(userId string, search string, sort string, order string) ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	setOrder := 1
	if order == "desc" {
		setOrder = -1
	}

	opts := options.Find()
	opts.SetSort(bson.D{
		{Key: sort, Value: setOrder},
	})

	escapedSearch := regexp.QuoteMeta(search)

	cursor, err := db.LinkCollection.Find(ctx, bson.M{"user_id": userId, "title": bson.M{"$regex": escapedSearch}}, opts)
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

func (LinkModel) GetAllLinkByUserIdAndLinkBookIdAndSearch(userId string, linkBookId string, search string, sort string, order string) ([]*Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var links []*Link

	// setOrder 를 만듭니다. asc, desc 에 따라서 정렬을 다르게 합니다.
	setOrder := 1
	if order == "desc" {
		setOrder = -1
	}

	opts := options.Find()
	opts.SetSort(bson.D{
		{Key: sort, Value: setOrder},
	})

	escapedSearch := regexp.QuoteMeta(search)

	cursor, err := db.LinkCollection.Find(ctx, bson.M{"user_id": userId, "link_book_id": linkBookId, "title": bson.M{"$regex": escapedSearch}}, opts)
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

func (LinkModel) DeleteAllLinksByUserId(userId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.LinkCollection.DeleteMany(ctx, bson.M{"user_id": userId})

	return result.DeletedCount, err
}

func (LinkModel) DeleteAllLinksByLinkIds(userId string, linkIds []string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{"link_id", bson.D{{"$in", linkIds}}},
		{"user_id", userId},
	}
	result, err := db.LinkCollection.DeleteMany(ctx, filter)

	return result.DeletedCount, err
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

func (LinkModel) UpdateTitleAndUrlByLinkId(linkId string, url string, title string, thumbnailURL string, tags []string) (*Link, error) {
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
		return nil, nil
	}

	/*
	 MongoDB의 공식 Go 언어 드라이버는 수정된 문서를 반환하는 기능을 제공하지 않습니다. 이는 MongoDB의 몇 가지 다른 드라이버와 달라 UpdateOne 후에 FindOne을 호출하는 것이 흔한 방법입니다.
	*/
	_, err := db.LinkCollection.UpdateOne(ctx, bson.M{"link_id": linkId}, bson.M{"$set": updateFields})

	if err != nil {
		return nil, err
	}

	link := &Link{}
	err = db.LinkCollection.FindOne(ctx, bson.M{"link_id": linkId}).Decode(link)

	if err != nil {
		return nil, err
	}

	return link, nil
}
