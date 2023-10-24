package notif

type NotificationUsecase struct {
	notificationModel NotificationModel
}

func (u NotificationUsecase) Notifications(userId string, page int64) (*NotificationRes, error) {
	result, err := u.notificationModel.Notifications(userId, page)
	if err != nil {
		return nil, err
	}
	return result, nil
}
