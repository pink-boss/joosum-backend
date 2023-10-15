package notif

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationUsecase struct {
	notificationModel NotificationModel
}

func (u NotificationUsecase) SaveDeviceId(deviceId, userId string) (*mongo.UpdateResult, error) {
	result, err := u.notificationModel.SaveDeviceId(deviceId, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u NotificationUsecase) GetNotificationAgree(userId string) (*mongo.UpdateResult, error) {
	//result, err := u.notificationModel.GetNotificationAgree(userId)
	//if err != nil {
	//	return nil, err
	//}
	//return result, nil
	return nil, nil
}

func (u NotificationUsecase) UpdatePushNotification(req PushNotificationReq, userId string) (*mongo.UpdateResult, error) {
	result, err := u.notificationModel.UpdatePushNotification(req, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
