package link

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"joosum-backend/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkBookModel struct{}

type LinkBookListReq struct {
	Sort string `form:"sort" enums:"created_at,last_saved_at,title"`
}

type LinkBookListRes struct {
	LinkBooks         []LinkBookRes `json:"linkBooks,omitempty"`
	TotalLinkCount    int64         `json:"TotalLinkCount,omitempty" example:"324"`
	NoFolderLinkCount int64         `json:"NoFolderLinkCount,omitempty" example:"13"`
}

type LinkBookRes struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"649028fab77fe1a8a3b0815e"`
	Title           string             `bson:"title" json:"title"`
	BackgroundColor string             `bson:"background_color" json:"backgroundColor"`
	TitleColor      string             `bson:"title_color" json:"titleColor"`
	Illustration    string             `bson:"illustration" json:"illustration"`
	CreatedAt       time.Time          `bson:"created_at"`
	LastSavedAt     time.Time          `bson:"last_saved_at"`
	UserId          string             `bson:"user_id" example:"User-0767d6af-a802-469c-9505-5ca91e03b354"`
	LinkCount       int64              `json:"linkCount"`
}

type LinkBookCreateReq struct {
	Title           string `json:"title"`
	BackgroundColor string `json:"backgroundColor"`
	TitleColor      string `json:"titleColor"`
	Illustration    string `json:"illustration"`
}

type LinkBook struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"649028fab77fe1a8a3b0815e"`
	Title           string             `bson:"title" json:"title"`
	BackgroundColor string             `bson:"background_color" json:"backgroundColor"`
	TitleColor      string             `bson:"title_color" json:"titleColor"`
	Illustration    string             `bson:"illustration" json:"illustration"`
	CreatedAt       time.Time          `bson:"created_at"`
	LastSavedAt     time.Time          `bson:"last_saved_at"`
	UserId          string             `bson:"user_id" example:"User-0767d6af-a802-469c-9505-5ca91e03b354"`
}

func (LinkBookModel) GetLinkBooks(req LinkBookListReq, userId string) ([]LinkBookRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 정렬 순서 디폴트: 생성 순
	if req.Sort == "" {
		req.Sort = "created_at"
	}

	// 폴더 정렬
	opts := options.Find()
	opts.SetSort(bson.M{req.Sort: -1})

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

	result, err := db.LinkBookCollection.InsertOne(ctx, linkBook)
	if err != nil {
		return nil, err
	}

	linkBook.ID = result.InsertedID.(primitive.ObjectID)

	return &linkBook, nil
}
