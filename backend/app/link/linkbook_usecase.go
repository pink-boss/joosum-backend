package link

import "time"

type LinkUsecase struct {
	linkModel LinkModel
}

func (u LinkUsecase) GetLinkBooks(req LinkBookListReq, userId string) ([]LinkBook, error) {
	res, err := u.linkModel.GetLinkBooks(req, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u LinkUsecase) CreateLinkBook(req LinkBookCreateReq, userId string) (interface{}, error) {

	link := LinkBook{
		Title:           req.Title,
		BackgroundColor: req.BackgroundColor,
		TitleColor:      req.TitleColor,
		Illustration:    req.Illustration,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		UserId:          userId,
	}

	res, err := u.linkModel.CreateLinkBook(link)
	if err != nil {
		return nil, err
	}
	return res, nil
}
