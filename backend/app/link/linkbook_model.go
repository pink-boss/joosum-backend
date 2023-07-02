package link

import (
	"context"
	"joosum-backend/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkBookModel struct{}

type LinkBookReq struct {
	Title           string `json:"title"`
	BackgroundColor string `json:"backgroundColor"`
	TitleColor      string `json:"titleColor"`
	Illustration    string `json:"illustration"`
}

type LinkBookRes struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"649028fab77fe1a8a3b0815e"`
	Title           string             `bson:"title" json:"title"`
	BackgroundColor string             `bson:"background_color" json:"backgroundColor"`
	TitleColor      string             `bson:"title_color" json:"titleColor"`
	Illustration    string             `bson:"illustration" json:"illustration"`
	CreatedAt       time.Time          `bson:"created_at" json:"-"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"-"`
	UserId          string             `bson:"user_id" example:"User-0767d6af-a802-469c-9505-5ca91e03b354"`
}

func (LinkBookModel) CreateLinkBook(link LinkBookRes) (*LinkBookRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.LinkCollection.InsertOne(ctx, link)
	if err != nil {
		return nil, err
	}

	link.ID = result.InsertedID.(primitive.ObjectID)

	return &link, nil
}
