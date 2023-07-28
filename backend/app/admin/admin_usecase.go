package admin

import (
	"joosum-backend/pkg/db"
)

type AdminUsecase struct {
	adminModel AdminModel
}

func (u AdminUsecase) UpdateDefaultFolder(req LinkBookUpdateReq, name string) (*db.LinkBook, error) {
	db.DefaultFolder.Title = req.Title
	db.DefaultFolder.TitleColor = req.TitleColor
	db.DefaultFolder.BackgroundColor = req.BackgroundColor
	db.DefaultFolder.Illustration = req.Illustration

	linkBook := db.LinkBook{
		Title:           req.Title,
		TitleColor:      req.TitleColor,
		BackgroundColor: req.BackgroundColor,
		Illustration:    req.Illustration,
		UpdatedBy:       name,
	}

	err := u.adminModel.UpdateDefaultFolder(linkBook)
	if err != nil {
		return nil, err
	}
	return &db.DefaultFolder, nil
}
