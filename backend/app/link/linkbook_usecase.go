package link

type LinkUsecase struct {
	linkModel LinkModel
}

func (u LinkUsecase) CreateLinkBook(req LinkBook) (interface{}, error) {
	link, err := u.linkModel.CreateLinkBook(req)
	if err != nil {
		return nil, err
	}
	return link, nil
}
