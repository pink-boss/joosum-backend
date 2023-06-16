package link

import (
	"context"
	"joosum-backend/pkg/db"
	"time"
)

type LinkModel struct{}

type LinkBook struct {
	Name         string `json:"name"`
	Color        string `json:"color"`
	Illustration string `json:"illustration"`
}

func (LinkModel) CreateLinkBook(req LinkBook) (*LinkBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.LinkCollection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
