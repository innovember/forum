package user

import (
	"github.com/innovember/forum/api/models"
)

type UserUsecase interface {
	Create(user *models.User) (status int, err error)
	GetAllUsers() (users []models.User, status int, err error)
}
