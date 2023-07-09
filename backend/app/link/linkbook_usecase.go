package link

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkBookUsecase struct {
	linkBookModel LinkBookModel
}

func (u LinkBookUsecase) GetLinkbooksForMainPage(userId string) ([]LinkBookRes, error) {

	linkBooks, err := u.linkBookModel.GetLinkBooks(LinkBookListReq{Sort: "last_saved_at"}, userId)
	if err != nil {
		return nil, err
	}

	// if linkBooks length 0 return []
	if len(linkBooks) == 0 {
		return []LinkBookRes{}, nil
	}

	return linkBooks, nil

}

func (u LinkBookUsecase) GetLinkBooks(req LinkBookListReq, userId string) (*LinkBookListRes, error) {
	linkBooks, err := u.linkBookModel.GetLinkBooks(req, userId)
	if err != nil {
		return nil, err
	}

	//todo total count 및 no folder count 추가

	//for _, linkBook := range linkBooks {
	//	linkCount, err := u.linkModel.GetLinkCount(linkBook.ID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	linkBook.LinkCount = *linkCount
	//}

	type Link struct {
		LinkId     string    `bson:"link_id" json:"linkId"`
		URL        string    `bson:"url" json:"url"`
		UserID     string    `bson:"user_id" json:"userId"`
		Title      string    `bson:"title" json:"title"`
		LinkBookId string    `bson:"link_book_id" json:"linkBookId"`
		ReadCount  int       `bson:"read_count" json:"readCount"`
		LastReadAt time.Time `bson:"last_read_at" json:"LastReadAt"`
		CreatedAt  time.Time `bson:"created_at" json:"CreatedAt"`
		UpdatedAt  time.Time `bson:"updated_at" json:"UpdatedAt"`
	}

	res := &LinkBookListRes{
		linkBooks,
		132,
	}

	return res, nil
}

func (u LinkBookUsecase) CreateLinkBook(req LinkBookCreateReq, userId string) (interface{}, error) {

	linkBook := LinkBook{
		ID:              uuid.New().String(),
		Title:           req.Title,
		BackgroundColor: req.BackgroundColor,
		TitleColor:      req.TitleColor,
		Illustration:    req.Illustration,
		CreatedAt:       time.Now(),
		UserId:          userId,
		IsDefault:       "n",
	}

	res, err := u.linkBookModel.CreateLinkBook(linkBook)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u LinkBookUsecase) CreateDefaultLinkBook(userId string) (interface{}, error) {

	linkBook := LinkBook{
		ID:              uuid.New().String(),
		Title:           "기본",
		BackgroundColor: "#6D6D6F",
		TitleColor:      "#FFFFFF",
		CreatedAt:       time.Now(),
		UserId:          userId,
		IsDefault:       "y",
	}

	res, err := u.linkBookModel.CreateLinkBook(linkBook)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u LinkBookUsecase) UpdateLinkBook(linkBookId string, req LinkBookCreateReq) (*mongo.UpdateResult, error) {

	linkBook := LinkBook{
		ID:              linkBookId,
		Title:           req.Title,
		BackgroundColor: req.BackgroundColor,
		TitleColor:      req.TitleColor,
		Illustration:    req.Illustration,
	}

	res, err := u.linkBookModel.UpdateLinkBook(linkBook)
	if err != nil {
		return nil, err
	}
	return res, nil
}
