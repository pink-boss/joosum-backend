package user

type UserUsecase struct {
	userModel UserModel
}

func (u UserUsecase)GetUserByEmail(email string) (*User, error) {
	user, err := u.userModel.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (u UserUsecase)RegisterUser(email string, socialType string) (*User, error) {
	// find email for check user exist
	isExistUser, err := u.GetUserByEmail(email)
	if err == nil {
		return nil, err
	}

	if isExistUser != nil {
		return nil, err
	}

	user, err := u.userModel.CreatUser(email, socialType)
	if err != nil {
		return nil, err
	}
	return user, err
}
