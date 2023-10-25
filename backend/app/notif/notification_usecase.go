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

func (u NotificationUsecase) ReadNotification(notificationId string) error {
	err := u.notificationModel.UpdateIsRead(notificationId)
	if err != nil {
		return err
	}
	return nil
}
