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

func (u UserUsecase) UpdateUserToInactiveUserByEmail(email string) error {

	user, err := u.userModel.FindUserByEmail(email)
	if err != nil {
		return err
	}

	err = u.userModel.DeleteUserByEmail(email)
	if err != nil {
		return err
	}

	err = u.userModel.CreateInactiveUserByUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u UserUsecase) GetWithdrawUsers() ([]*InactiveUser, error) {
	users, err := u.userModel.FindInactiveUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserUsecase) GetInactiveUserByEmail(email string) (*InactiveUser, error) {
	inactiveUser, err := u.userModel.FindInactiveUser(email)
	if err != nil {
		return nil, err
	}
	return inactiveUser, err
}
