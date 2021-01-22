package user

import (
	"github.com/innovember/forum/api/models"
)

type UserRepository interface {
	Create(user *models.User) (status int, err error)
	CheckByUsernameOrEmail(user *models.User) (status int, err error)
}
