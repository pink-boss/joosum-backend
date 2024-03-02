package tag

import (
	"joosum-backend/pkg/util"
)

type TagUsecase struct {
	tagModel TagModel
}

func (u TagUsecase) CreateTags(userId string, names []string) (*Tag, error) {
	tag, err := u.tagModel.UpsertTags(userId, names)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (u TagUsecase) FindTagsByUserId(userId string) ([]string, error) {
	tags, err := u.tagModel.FindTagsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (u TagUsecase) DeleteTag(userId string, tag string) ([]string, error) {
	tags, err := u.tagModel.FindTagsByUserId(userId)
	if err != nil {
		return nil, err
	}
	tags = util.ListUtil.Remove(tags, tag)
	_, err = u.tagModel.UpsertTags(userId, tags)
	if err != nil {
		return nil, err
	}

	return tags, err
}
