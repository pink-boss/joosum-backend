package tag

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

func (u TagUsecase) DeleteTag(user_id string, tag_id string) error {
	err := u.tagModel.DeleteTag(user_id, tag_id)
	if err != nil {
		return err
	}
	return nil
}
