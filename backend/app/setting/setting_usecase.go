package setting

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type SettingUsecase struct {
	settingModel SettingModel
}

func (u SettingUsecase) SaveDeviceId(deviceId, userId string) (*mongo.UpdateResult, error) {
	result, err := u.settingModel.SaveDeviceId(deviceId, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u SettingUsecase) GetNotificationAgree(userId string) (*NotificationAgree, error) {
	result, err := u.settingModel.GetNotificationAgree(userId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			agree := &NotificationAgree{
				IsReadAgree:     false,
				IsClassifyAgree: false,
				UserId:          userId,
			}
			return agree, nil
		}
		return nil, err
	}
	return result, nil
}

func (u SettingUsecase) UpdatePushNotification(req PushNotificationReq, userId string) (*mongo.UpdateResult, error) {
	result, err := u.settingModel.UpdatePushNotification(req, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
