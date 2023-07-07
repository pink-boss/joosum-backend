package link

import "time"

type LinkBookUsecase struct {
	linkBookModel LinkBookModel
}

func (u LinkBookUsecase) GetLinkBooks(req LinkBookListReq, userId string) (*LinkBookListRes, error) {
	linkBooks, err := u.linkBookModel.GetLinkBooks(req, userId)
	if err != nil {
		return nil, err
	}

	// todo total count 및 no folder count 추가

	res := &LinkBookListRes{
		linkBooks,
		132,
	}

	return res, nil
}

func (u LinkBookUsecase) CreateLinkBook(req LinkBookCreateReq, userId string) (interface{}, error) {

	linkBook := LinkBook{
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
