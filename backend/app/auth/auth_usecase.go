package auth

import (
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
)

type AuthUsecase struct {
	salt      string
	userModel user.UserModel
}

func (u *AuthUsecase) GenerateNewJWTToken(roles []string, email string) (string, string, error) {
	accessToken, err := util.GenerateNewJWTAccessToken([]string{"USER", "ADMIN"}, email)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := util.GenerateNewJWTRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *AuthUsecase) SignUp(userInfo user.User) (*user.User, error) {
	user, err := u.userModel.CreatUser(userInfo)
	if err != nil {
		return nil, err
	}
	return user, nil
}
