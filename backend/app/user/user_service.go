package user

func GetUserByEmail(email string) (*User, error) {
	user, err := FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, err
}
