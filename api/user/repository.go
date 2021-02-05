package user

import (
	"github.com/innovember/forum/api/models"
)

type UserRepository interface {
	Create(user *models.User) (status int, err error)
	CheckByUsernameOrEmail(user *models.User) (status int, err error)
	GetAllUsers() (users []models.User, err error)
	GetUserByID(userID int64) (user *models.User, err error)
	GetPassword(username string) (password string, status int, err error)
	FindUserByUsername(username string) (user *models.User, status int, err error)
	UpdateSession(userID int64, sessionValue string) (err error)
	ValidateSession(sessionValue string) (user *models.User, status int, err error)
	CheckSessionByUsername(username string) (status int, err error)
}
