package db

import "time"

var DefaultFolder LinkBook

type LinkBook struct {
	Title           string    `bson:"title" example:"기본"`
	TitleColor      string    `bson:"title_color" example:"#FFFFFF"`
	BackgroundColor string    `bson:"background_color" example:"#8A8A9A"`
	Illustration    *string   `bson:"illustration" example:"illust11"`
	UpdatedAt       time.Time `bson:"updated_at" swaggerignore:"true"`
	UpdatedBy       string    `bson:"updated_by" swaggerignore:"true"`
}
