package link

import (
	"context"
	"joosum-backend/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/errgo.v2/errors"
)

type LinkBookModel struct{}

type LinkBookListReq struct {
	Sort string `form:"sort" enums:"created_at,last_saved_at,title"`
}

type LinkBookListRes struct {
	LinkBooks      []LinkBookRes `json:"linkBooks"`
	TotalLinkCount int64         `json:"totalLinkCount" example:"324"`
}

type LinkBookRes struct {
	LinkBookId      string    `bson:"_id" json:"linkBookId" example:"649028fab77fe1a8a3b0815e"`
	Title           string    `bson:"title" json:"title"`
	BackgroundColor string    `bson:"background_color" json:"backgroundColor"`
	TitleColor      string    `bson:"title_color" json:"titleColor"`
	Illustration    *string   `bson:"illustration" json:"illustration"`
	CreatedAt       time.Time `bson:"created_at" json:"createdAt"`
	LastSavedAt     time.Time `bson:"last_saved_at" json:"lastSavedAt"`
	UserId          string    `bson:"user_id" example:"User-0767d6af-a802-469c-9505-5ca91e03b354" json:"userId"`
	LinkCount       int64     `json:"linkCount"`
	IsDefault       string    `bson:"is_default" json:"isDefault"`
}

type LinkBookCreateReq struct {
	Title           string  `json:"title" example:"title" validate:"required"`
	BackgroundColor string  `json:"backgroundColor" example:"#6D6D6F"`
	TitleColor      string  `json:"titleColor" example:"#FFFFFF"`
	Illustration    *string `json:"illustration"`
}

type LinkBookDeleteRes struct {
	DeletedLinks int64 `json:"deletedLinks"`
}

type LinkBook struct {
	LinkBookId      string    `bson:"_id" json:"linkBookId" example:"649028fab77fe1a8a3b0815e"`
	Title           string    `bson:"title" json:"title"`
	BackgroundColor string    `bson:"background_color" json:"backgroundColor"`
	TitleColor      string    `bson:"title_color" json:"titleColor"`
	Illustration    *string   `bson:"illustration" json:"illustration"`
	CreatedAt       time.Time `bson:"created_at" json:"createdAt"`
	LastSavedAt     time.Time `bson:"last_saved_at" json:"lastSavedAt"`
	UserId          string    `bson:"user_id" example:"User-0767d6af-a802-469c-9505-5ca91e03b354" json:"userId"`
	IsDefault       string    `bson:"is_default" json:"isDefault"`
}

func (LinkBookModel) GetLinkBooks(req LinkBookListReq, userId string) ([]LinkBookRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 정렬 순서 디폴트: 생성 순
	if req.Sort == "" || req.Sort == "create_at" {
		req.Sort = "created_at"
	}

	var order int
	switch req.Sort {
	case "created_at":
		order = -1 // db.Desc와 동일
	case "title": // 폴더명 순일 때는 오름차순
		order = 1 // db.Asc와 동일
	case "last_saved_at":
		order = -1 // db.Desc와 동일
	default:
		order = -1 // 기본값은 내림차순
	}

	sort := bson.D{
		{"is_default", -1}, // 기본 폴더북은 정렬에 관계없이 첫번째. y-n 정렬 (lmn opqr...vwxyz)
		{req.Sort, order},
	}

	if req.Sort == "last_saved_at" {
		sort = append(sort, bson.E{"created_at", -1}) // 업데이트 순이 같다면 생성 순으로 정렬
	}
	// 폴더 정렬
	opts := options.Find()
	opts.SetSort(sort)
	cur, err := db.LinkBookCollection.Find(ctx, map[string]string{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var linkBooks []LinkBookRes

	for cur.Next(ctx) {
		var result LinkBookRes
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}

		linkBooks = append(linkBooks, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return linkBooks, nil
}

func (LinkBookModel) CreateLinkBook(linkBook LinkBook) (*LinkBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.LinkBookCollection.InsertOne(ctx, linkBook)
	if err != nil {
		return nil, err
	}

	return &linkBook, nil
}

func (LinkBookModel) GetLinkBookById(linkBookId string) (*LinkBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var linkBook *LinkBook
	error := db.LinkBookCollection.FindOne(ctx, bson.M{
		"_id": linkBookId,
	}).Decode(&linkBook)

	if error != nil {
		return nil, error
	}

	return linkBook, nil
}

func (LinkBookModel) GetDefaultLinkBook(userId string) (*LinkBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var linkBook *LinkBook
	error := db.LinkBookCollection.FindOne(ctx, bson.M{
		"user_id":    userId,
		"is_default": "y",
	}).Decode(&linkBook)

	if error != nil {
		return nil, error
	}

	return linkBook, nil
}

func (LinkBookModel) UpdateLinkBook(linkBook LinkBook) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	linkBookUpdate := bson.M{
		"$set": bson.M{
			"title":            linkBook.Title,
			"background_color": linkBook.BackgroundColor,
			"title_color":      linkBook.TitleColor,
			"illustration":     linkBook.Illustration,
		},
	}

	linkUpdate := bson.M{
		"$set": bson.M{
			"link_book_name": linkBook.Title,
		},
	}

	err := db.LinkBookCollection.FindOneAndUpdate(ctx, bson.M{"_id": linkBook.LinkBookId}, linkBookUpdate).Decode(&mongo.SingleResult{})
	if err == mongo.ErrNoDocuments {
		return errors.New(err.Error())
	}

	_, err = db.LinkCollection.UpdateMany(ctx, bson.M{"link_book_id": linkBook.LinkBookId}, linkUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (LinkBookModel) UpdateLinkBookLastSavedAt(linkBookId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"last_saved_at": time.Now(),
		},
	}

	result := db.LinkBookCollection.FindOneAndUpdate(ctx, bson.M{"_id": linkBookId}, update).Decode(&mongo.SingleResult{})
	if result == mongo.ErrNoDocuments {
		return errors.New(result.Error())
	}

	return result
}

func (LinkBookModel) DeleteLinkBook(linkBookId string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.LinkBookCollection.DeleteOne(ctx, bson.M{"_id": linkBookId})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (LinkBookModel) IsDefaultLinkBook(linkBookId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := db.LinkBookCollection.CountDocuments(ctx, bson.M{"_id": linkBookId, "is_default": "y"})
	if err != nil {
		return false, err
	}

	if count == 1 {
		return true, err
	}
	return false, nil
}

func (LinkBookModel) IsDuplicatedTitle(title, userId string, linkBookId *string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := db.LinkBookCollection.CountDocuments(ctx, bson.M{"title": title, "user_id": userId, "_id": bson.M{"$ne": linkBookId}})
	if err != nil {
		return false, err
	}

	if count >= 1 {
		return true, err
	}
	return false, nil
}

func (LinkBookModel) HaveLinkBook(linkBookId string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := db.LinkBookCollection.FindOne(ctx, bson.M{"_id": linkBookId}).Decode(&mongo.SingleResult{})
	if result == mongo.ErrNoDocuments {
		return false
	}
	return true
}
