package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
)

type NotificationUsecase struct {
	notificationRepo post.NotificationRepository
}

func NewNotificationUsecase(repo post.NotificationRepository) post.NotificationUsecase {
	return &NotificationUsecase{notificationRepo: repo}
}

func (nu *NotificationUsecase) Create(notification *models.Notification) (newNotification *models.Notification, status int, err error) {
	if newNotification, status, err = nu.notificationRepo.Create(notification); err != nil {
		return nil, status, err
	}
	return newNotification, status, err
}
func (nu *NotificationUsecase) DeleteAllNotifications(receiverID int64) (status int, err error) {
	if status, err = nu.notificationRepo.DeleteAllNotifications(receiverID); err != nil {
		return status, err
	}
	return status, err
}
func (nu *NotificationUsecase) GetAllNotifications(receiverID int64) (notifications []models.Notification, status int, err error) {
	if notifications, status, err = nu.notificationRepo.GetAllNotifications(receiverID); err != nil {
		return nil, status, err
	}
	return notifications, status, nil
}
