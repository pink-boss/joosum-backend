package user

type UserUsecase struct {
	userModel UserModel
}

func (u UserUsecase) GetUserByEmail(email string) (*User, error) {
	user, err := u.userModel.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (u UserUsecase) GetUserById(id string) (*User, error) {
	user, err := u.userModel.FindUserById(id)
	if err != nil {
		return nil, err
	}
	return user, err
}
