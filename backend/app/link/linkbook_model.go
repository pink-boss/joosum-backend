package link

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"joosum-backend/pkg/db"
	"time"
)

type LinkBookModel struct{}

type LinkBookListReq struct {
	Sort string `form:"sort" enums:"created_at,last_saved_at,title"`
}

type LinkBookListRes struct {
	LinkBooks      []LinkBookRes `json:"linkBooks,omitempty"`
	TotalLinkCount int64         `json:"totalLinkCount,omitempty" example:"324"`
}

type LinkBookRes struct {
	LinkBookId      string    `bson:"_id,omitempty" json:"linkBookId" example:"649028fab77fe1a8a3b0815e"`
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
	LinkBookId      string    `bson:"_id,omitempty" json:"linkBookId" example:"649028fab77fe1a8a3b0815e"`
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

	sort := -1

	// 정렬 순서 디폴트: 생성 순
	if req.Sort == "" {
		req.Sort = "created_at"
	}

	// 폴더명 순일 때는 오름차순
	if req.Sort == "title" {
		sort = 1
	}

	// 폴더 정렬
	opts := options.Find()
	opts.SetSort(bson.D{
		{Key: "is_default", Value: 1}, // ""-n-y 정렬 (lmn opqr...vwxyz)
		{Key: req.Sort, Value: sort},
	})

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

	// todo 링크북에 링크 수 카운트 추가

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

func (LinkBookModel) UpdateLinkBook(linkBook LinkBook) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":            linkBook.Title,
			"background_color": linkBook.BackgroundColor,
			"title_color":      linkBook.TitleColor,
			"illustration":     linkBook.Illustration,
		},
	}

	result, err := db.LinkBookCollection.UpdateByID(ctx, linkBook.LinkBookId, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (LinkBookModel) UpdateLinkBookLastSavedAt(linkBookId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"last_saved_at": time.Now(),
		},
	}

	_, err := db.LinkBookCollection.UpdateByID(ctx, linkBookId, update)
	if err != nil {
		return err
	}

	return nil
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
