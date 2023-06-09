package link

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/google/uuid"
)

type LinkBookUsecase struct {
	linkModel     LinkModel
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

	for i, linkBook := range linkBooks {
		linkCount, err := u.linkModel.GetLinkBookLinkCount(linkBook.LinkBookId)
		if err != nil {
			return nil, err
		}
		linkBooks[i].LinkCount = linkCount
	}

	totalLinkCount, err := u.linkModel.GetUserLinkCount(userId)
	if err != nil {
		return nil, err
	}

	res := &LinkBookListRes{
		linkBooks,
		totalLinkCount,
	}

	return res, nil
}

func (u LinkBookUsecase) CreateLinkBook(req LinkBookCreateReq, userId string) (interface{}, error) {

	linkBook := LinkBook{
		LinkBookId:      "LinkBook-" + uuid.New().String(),
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
		LinkBookId:      "LinkBook-" + uuid.New().String(),
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
		LinkBookId:      linkBookId,
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

func (u LinkBookUsecase) DeleteLinkBookWithLinks(userId, linkBookId string) (*LinkBookDeleteRes, error) {
	isDefault, err := u.linkBookModel.IsDefaultLinkBook(linkBookId)
	if err != nil {
		return nil, err
	}

	// 기본 링크북 폴더는 삭제하지 않음
	if !isDefault {
		_, err := u.linkBookModel.DeleteLinkBook(linkBookId)
		if err != nil {
			return nil, err
		}
	}

	result, err := u.linkModel.DeleteAllLinksByLinkBookId(userId, linkBookId)
	if err != nil {
		return nil, err
	}

	return &LinkBookDeleteRes{DeletedLinks: result.DeletedCount}, nil
}
