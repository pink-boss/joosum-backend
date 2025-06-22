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
	tagData, err := u.tagModel.FindTagByUserId(userId)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)

	used := make(map[string]bool)
	for _, tag := range tagData.LastUsed {
		for _, name := range tagData.Names {
			if tag == name && !used[tag] {
				tags = append(tags, tag)
				used[tag] = true
			}
		}
	}

	for _, name := range tagData.Names {
		if !used[name] {
			tags = append(tags, name)
			used[name] = true
		}
	}

	return tags, nil
}

// FindTagsByUserIdAndSearch는 사용자 아이디와 검색어로 태그를 조회합니다.
func (u TagUsecase) FindTagsByUserIdAndSearch(userId string, search string) ([]string, error) {
	tagData, err := u.tagModel.FindTagByUserId(userId)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	used := make(map[string]bool)

	for _, tag := range tagData.LastUsed {
		for _, name := range tagData.Names {
			if tag == name && util.HangulMatch(tag, search) && !used[tag] {
				tags = append(tags, tag)
				used[tag] = true
			}
		}
	}

	for _, name := range tagData.Names {
		if util.HangulMatch(name, search) && !used[name] {
			tags = append(tags, name)
			used[name] = true
		}
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

// UpdateLastUsedTags는 사용자가 링크를 저장할 때 사용한 태그 순서를 기록합니다.
func (u TagUsecase) UpdateLastUsedTags(userId string, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	return u.tagModel.UpdateLastUsedTags(userId, tags)
}
