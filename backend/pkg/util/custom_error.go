package util

import "gopkg.in/errgo.v2/errors"

var (
	ErrDuplicatedTitle = errors.New("같은 이름의 폴더가 존재합니다")
	ErrInvalidSort     = errors.New("허용되지 않은 sort 값입니다")
)
