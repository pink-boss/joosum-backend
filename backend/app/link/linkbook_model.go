package link

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"joosum-backend/pkg/db"
	"log"
	"time"
)

type LinkModel struct{}

type LinkBookListReq struct {
	Title     string `json:"title"`
	CreatedAt int8   `json:"createdAt"`
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
	UpdatedAt       time.Time          `bson:"updated_at"`
	UserId          string             `bson:"user_id" example:"User-0767d6af-a802-469c-9505-5ca91e03b354"`
}

func (LinkModel) GetLinkBooks(req LinkBookListReq, userId string) ([]LinkBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.M{"created_at": -1})

	cur, err := db.LinkBookCollection.Find(ctx, map[string]string{"user_id": userId}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var linkBooks []LinkBook

	for cur.Next(ctx) {
		var result LinkBook
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		linkBooks = append(linkBooks, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return linkBooks, nil
}

func (LinkModel) CreateLinkBook(link LinkBook) (*LinkBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.LinkBookCollection.InsertOne(ctx, link)
	if err != nil {
		return nil, err
	}

	link.ID = result.InsertedID.(primitive.ObjectID)

	return &link, nil
}
