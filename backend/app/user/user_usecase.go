package user

func GetUserByEmail(email string) (*User, error) {
	user, err := FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func RegisterUser(email string, socialType string) (*User, error) {
	// find email for check user exist
	isExistUser, err := GetUserByEmail(email)
	if err == nil {
		return nil, err
	}

	if isExistUser != nil {
		return nil, err
	}

	user, err := CreatUser(email, socialType)
	if err != nil {
		return nil, err
	}
	return user, err
}
