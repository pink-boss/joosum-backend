package tag

type TagUsecase struct {
	tagModel TagModel
}

func (u TagUsecase) CreateTag(name string, user_id string) (*Tag, error) {
	tag, err := u.tagModel.CreateTag(name, user_id)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (u TagUsecase) FindTagByUserId(user_id string) ([]*Tag, error) {
	tags, err := u.tagModel.FindTagByUserId(user_id)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (u TagUsecase) DeleteTag(user_id string,tag_id string)  error {
	err := u.tagModel.DeleteTag(user_id, tag_id)
	if err != nil {
		return err
	}
	return nil
}