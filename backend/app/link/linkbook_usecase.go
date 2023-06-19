package link

import "time"

type LinkUsecase struct {
	linkModel LinkModel
}

func (u LinkUsecase) CreateLinkBook(req LinkBookReq, userId string) (interface{}, error) {

	link := LinkBookRes{
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
