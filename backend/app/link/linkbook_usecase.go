package link

import "time"

type LinkUsecase struct {
	linkModel LinkModel
}

func (u LinkUsecase) GetLinkBooks(req LinkBookListReq, userId string) (*LinkBookListRes, error) {
	linkBooks, err := u.linkModel.GetLinkBooks(req, userId)
	if err != nil {
		return nil, err
	}

	// todo total count 및 no folder count 추가

	res := &LinkBookListRes{
		linkBooks,
		132,
		13,
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
		UserId:          userId,
	}

	res, err := u.linkModel.CreateLinkBook(link)
	if err != nil {
		return nil, err
	}
	return res, nil
}
