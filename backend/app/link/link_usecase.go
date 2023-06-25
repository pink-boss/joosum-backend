package link

type LinkUsecase struct {
	linkModel LinkModel
}

func (u LinkUsecase) CreateLink(url string, title string, userId string, linkBookId string) (*Link, error) {
	link, err := u.linkModel.CreateLink(url, title, userId, linkBookId)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (u LinkUsecase) FindOneLinkByLinkId(linkId string) (*Link, error) {
	link, err := u.linkModel.GetOneLinkByLinkId(linkId)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (u LinkUsecase) FindAllLinksByUserId(userId string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserId(userId)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (u LinkUsecase) FindAllLinksByUserIdAndLinkBookId(userId string, linkBookId string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndLinkBookId(userId, linkBookId)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (u LinkUsecase) DeleteOneByLinkId(linkId string) error {
	err := u.linkModel.DeleteOneByLinkId(linkId)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) DeleteAllLinksByUserId(userId string) error {
	err := u.linkModel.DeleteAllLinksByUserId(userId)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) DeleteAllLinksByLinkBookId(userId string, linkBookId string) error {
	err := u.linkModel.DeleteAllLinksByLinkBookId(userId, linkBookId)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) UpdateReadByLinkId(linkId string) error {
	err := u.linkModel.UpdateReadCountByLinkId(linkId)
	if err != nil {
		return err
	}

	return nil

}

func (u LinkUsecase) UpdateLinkBookIdByLinkId(linkId string, linkBookId string) error {
	err := u.linkModel.UpdateLinkBookIdByLinkId(linkId, linkBookId)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) UpdateTitleAndUrlByLinkId(linkId string, url string, title string) error {
	err := u.linkModel.UpdateTitleAndUrlByLinkId(linkId, url, title)
	if err != nil {
		return err
	}

	return nil
}
