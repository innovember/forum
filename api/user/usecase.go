package user

import (
	"github.com/innovember/forum/api/models"
)

type UserUsecase interface {
	Create(user *models.User) (status int, err error)
	GetAllUsers() (users []models.User, status int, err error)
	GetUserByID(userID int64) (user *models.User, err error)
	GetPassword(username string) (password string, status int, err error)
	FindUserByUsername(username string) (user *models.User, status int, err error)
	UpdateSession(userID int64, sessionValue string, expiresAt int64) (err error)
	ValidateSession(sessionValue string) (user *models.User, status int, err error)
	CheckSessionByUsername(username string) (status int, err error)
	UpdateActivity(userID int64) (err error)
	CreateRoleRequest(userID int64) (err error)
	GetRoleRequestByUserID(userID int64) (request *models.RoleRequest, err error)
	DeleteRoleRequest(userID int64) (err error)
}

type AdminUsecase interface {
	UpgradeRole(userID int64) (err error)
	GetAllRoleRequests() (roleRequests []models.RoleRequest, err error)
	DeleteRoleRequest(userID int64) (err error)
}
