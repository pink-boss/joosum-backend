package admin

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/errgo.v2/errors"
	"joosum-backend/pkg/db"
	"time"
)

type AdminModel struct{}

type LinkBookUpdateReq struct {
	Title           string  `json:"title" example:"기본" validate:"required"`
	BackgroundColor string  `json:"backgroundColor" example:"#8A8A9A" validate:"required"`
	TitleColor      string  `json:"titleColor" example:"#FFFFFF" validate:"required"`
	Illustration    *string `json:"illustration" example:"illust11"`
}

func (AdminModel) UpdateDefaultFolder(linkBook db.LinkBook) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":            linkBook.Title,
			"background_color": linkBook.BackgroundColor,
			"title_color":      linkBook.TitleColor,
			"illustration":     linkBook.Illustration,
			"updated_by":       linkBook.UpdatedBy,
			"updated_at":       time.Now(),
		},
	}

	err := db.CommonCollection.FindOneAndUpdate(ctx, bson.M{"type": "DEFAULT_FOLDER"}, update).Decode(&mongo.SingleResult{})
	if err == mongo.ErrNoDocuments {
		return errors.New(err.Error())
	}
	return nil
}
