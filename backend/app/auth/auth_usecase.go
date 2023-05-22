package auth

import (
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
)

type AuthUsecae struct {
	salt string
}

func (u *AuthUsecae) GenerateNewJWTToken(roles []string, email string) (string, string, error) {
	accessToken, err := util.GenerateNewJWTAccessToken([]string{"USER", "ADMIN"}, email, u.salt)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := util.GenerateNewJWTRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *AuthUsecae) SignUp(email string, social string) (*user.User, error) {
	user, err := user.CreatUser(email, social)
	if err != nil {
		return nil, err
	}
	return user, nil
}
