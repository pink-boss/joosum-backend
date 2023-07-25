package util

import "gopkg.in/errgo.v2/errors"

var ErrDuplicatedTitle = errors.New("같은 이름의 폴더가 존재합니다")
