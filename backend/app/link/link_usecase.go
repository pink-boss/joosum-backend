package link

import "joosum-backend/pkg/db"

type LinkUsecase struct {
	linkModel     LinkModel
	linkBookModel LinkBookModel
}

func (u LinkUsecase) CreateLink(url string, title string, userId string, linkBookId string, thumbnailURL string, tags []string) (*Link, error) {

	// linkBookId 가 root 이거나 빈 스트링이라면 기본 폴더에 저장
	var linkBookName string
	if linkBookId == "root" || linkBookId == "" {
		defaultLinkBook, err := u.linkBookModel.GetDefaultLinkBook(userId)
		if err != nil {
			return nil, err
		}

		linkBookId = defaultLinkBook.LinkBookId
		linkBookName = db.DefaultFolder.Title
	} else {
		linkBookData, err := u.linkBookModel.GetLinkBookById(linkBookId)
		if err != nil {
			return nil, err
		}

		linkBookName = linkBookData.Title
	}

	link, err := u.linkModel.CreateLink(url, title, userId, linkBookId, linkBookName, thumbnailURL, tags)
	if err != nil {
		return nil, err
	}

	// 링크북 최근 링크등록일 업데이트
	err = u.linkBookModel.UpdateLinkBookLastSavedAt(linkBookId)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (u LinkUsecase) Get9LinksByUserId(userId string) ([]*Link, error) {
	links, err := u.linkModel.Get9LinksByUserId(userId)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (u LinkUsecase) FindOneLinkByLinkId(linkId string) (*Link, error) {
	link, err := u.linkModel.GetOneLinkByLinkId(linkId)
	if err != nil {
		return nil, err
	}

	isDefault, err := u.linkBookModel.IsDefaultLinkBook(link.LinkBookId)
	if err != nil {
		return nil, err
	}
	if isDefault {
		link.LinkBookName = db.DefaultFolder.Title
	}

	return link, nil
}

func (u LinkUsecase) FindAllLinksByUserId(userId string, sort string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserId(userId, sort)
	if err != nil {
		return nil, err
	}

	// if links length 0 return []
	if len(links) == 0 {
		return []*Link{}, nil
	}

	for _, link := range links {
		isDefault, err := u.linkBookModel.IsDefaultLinkBook(link.LinkBookId)
		if err != nil {
			return nil, err
		}
		if isDefault {
			link.LinkBookName = db.DefaultFolder.Title
		}
	}

	return links, nil
}

func (u LinkUsecase) FindAllLinksByUserIdAndSearch(userId string, search string, sort string, order string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndSearch(userId, search, sort, order)
	if err != nil {
		return nil, err
	}

	// if links length 0 return []
	if len(links) == 0 {
		return []*Link{}, nil
	}

	for _, link := range links {
		isDefault, err := u.linkBookModel.IsDefaultLinkBook(link.LinkBookId)
		if err != nil {
			return nil, err
		}
		if isDefault {
			link.LinkBookName = db.DefaultFolder.Title
		}
	}

	return links, nil
}

func (u LinkUsecase) FindAllLinksByUserIdAndLinkBookId(userId string, linkBookId string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndLinkBookId(userId, linkBookId)
	if err != nil {
		return nil, err
	}

	isDefault, err := u.linkBookModel.IsDefaultLinkBook(linkBookId)
	if isDefault {
		for _, link := range links {
			link.LinkBookName = db.DefaultFolder.Title
		}
	}

	return links, nil
}

func (u LinkUsecase) FindAllLinksByUserIdAndLinkBookIdAndSearch(userId string, linkBookId string, search string, sort string, order string) ([]*Link, error) {
	links, err := u.linkModel.GetAllLinkByUserIdAndLinkBookIdAndSearch(userId, linkBookId, search, sort, order)
	if err != nil {
		return nil, err
	}

	// if links length 0 return []
	if len(links) == 0 {
		return []*Link{}, nil
	}

	isDefault, err := u.linkBookModel.IsDefaultLinkBook(linkBookId)
	if isDefault {
		for _, link := range links {
			link.LinkBookName = db.DefaultFolder.Title
		}
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
	_, err := u.linkModel.DeleteAllLinksByLinkBookId(userId, linkBookId)
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
	// find LinkBook by linkBookId
	linkBookData, err := u.linkBookModel.GetLinkBookById(linkBookId)
	if err != nil {
		return err
	}

	linkBookName := linkBookData.Title

	err = u.linkModel.UpdateLinkBookIdAndTitleByLinkId(linkId, linkBookId, linkBookName)
	if err != nil {
		return err
	}

	return nil
}

func (u LinkUsecase) UpdateTitleAndUrlByLinkId(linkId string, url string, title string, thumbnailURL string, tags []string) (*Link, error) {
	link, err := u.linkModel.UpdateTitleAndUrlByLinkId(linkId, url, title, thumbnailURL, tags)
	if err != nil {
		return nil, err
	}

	return link, nil
}
