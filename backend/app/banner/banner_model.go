package banner

import (
	"bytes"
	"context"
	"joosum-backend/pkg/db"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type BannerModel struct {
}

type Banner struct {
	ID       string `json:"id"`
	ImageURL string `json:"imageURL"`
	ClickURL string `json:"clickURL"`
}

type BannerCreateReq struct {
	ImageURL string `json:"imageURL" example:"https://example.com/image.jpg"`
	ClickURL string `json:"clickURL" example:"https://example.com"`
}

func (BannerModel) GetBanners() ([]*Banner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var banners []*Banner

	cursor, error := db.BannerCollection.Find(ctx, bson.M{})
	if error != nil {
		return nil, error
	}

	for cursor.Next(ctx) {
		var banner Banner
		cursor.Decode(&banner)
		banners = append(banners, &banner)
	}

	return banners, nil

}

func (BannerModel) CreateBanner(imageURL string, clickURL string) (*Banner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uniqueKey := make(chan string)

	go func() {
		s := "Banner-"
		buf := bytes.NewBufferString(s)
		uid := uuid.New()
		buf.WriteString(uid.String())
		uniqueKey <- buf.String()
	}()

	uid := <-uniqueKey

	banner := Banner{
		ID:       uid,
		ImageURL: imageURL,
		ClickURL: clickURL,
	}

	_, error := db.BannerCollection.InsertOne(ctx, banner)
	if error != nil {
		return nil, error
	}

	return &banner, nil
}

func (BannerModel) DeleteBannerById(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, error := db.BannerCollection.DeleteOne(ctx, bson.M{"id": id})
	if error != nil {
		return 0, error
	}

	return result.DeletedCount, nil
}
