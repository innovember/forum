package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/user"
)

type UserNotificationUsecase struct {
	userNotificationRepo user.UserNotificationRepository
}

func NewUserNotificationUsecase(repo user.UserNotificationRepository) user.UserNotificationUsecase {
	return &UserNotificationUsecase{userNotificationRepo: repo}
}

func (uu *UserNotificationUsecase) CreateRoleNotification(roleNotification *models.RoleNotification) (err error) {
	if err = uu.userNotificationRepo.CreateRoleNotification(roleNotification); err != nil {
		return err
	}
	return nil
}
func (uu *UserNotificationUsecase) CreatePostReportNotification(postNotification *models.PostReportNotification) (err error) {
	if err = uu.userNotificationRepo.CreatePostReportNotification(postNotification); err != nil {
		return err
	}
	return nil
}

func (uu *UserNotificationUsecase) DeleteAllRoleNotifications(userID int64) (err error) {
	if err = uu.userNotificationRepo.DeleteAllRoleNotifications(userID); err != nil {
		return err
	}
	return nil
}

func (uu *UserNotificationUsecase) DeleteAllPostReportNotifications(userID int64) (err error) {
	if err = uu.userNotificationRepo.DeleteAllPostReportNotifications(userID); err != nil {
		return err
	}
	return nil
}

func (uu *UserNotificationUsecase) GetRoleNotifications(userID int64) (roleNotifications []models.RoleNotification, err error) {
	if roleNotifications, err = uu.userNotificationRepo.GetRoleNotifications(userID); err != nil {
		return nil, err
	}
	return roleNotifications, nil
}

func (uu *UserNotificationUsecase) GetPostReportNotifications(userID int64) (postNotifications []models.PostReportNotification, err error) {
	if postNotifications, err = uu.userNotificationRepo.GetPostReportNotifications(userID); err != nil {
		return nil, err
	}
	return postNotifications, nil
}
