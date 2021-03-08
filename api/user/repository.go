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
	UpdateSession(userID int64, sessionValue string, expiresAt int64) (err error)
	ValidateSession(sessionValue string) (user *models.User, status int, err error)
	CheckSessionByUsername(username string) (status int, err error)
	UpdateActivity(userID int64) (err error)
	CreateRoleRequest(userID int64) (err error)
	GetRoleRequestByUserID(userID int64) (request *models.RoleRequest, err error)
	DeleteRoleRequest(userID int64) (err error)
}

type AdminRepository interface {
	UpgradeRole(userID int64) (err error)
	GetAllRoleRequests() (roleRequests []models.RoleRequest, err error)
	DeleteRoleRequest(userID int64) (err error)
	GetAllPostReports() (postReports []models.PostReport, err error)
	AcceptPostReport(postReportID int64) (err error)
	DismissPostReport(postReportID int64) (err error)
	GetAllModerators() (moderators []models.User, err error)
	DemoteModerator(moderatorID int64) (err error)
}

type ModeratorRepository interface {
	CreatePostReport(postReport *models.PostReport) (err error)
	DeletePostReport(postReportID int64) (err error)
	GetMyReports(moderatorID int64) (postReports []models.PostReport, err error)
	ApprovePost(postID int64) (err error)
	GetAllUnapprovedPosts() (posts []models.Post, err error)
	BanPost(postID int64, bans []string) (err error)
}

type UserNotificationRepository interface {
	CreateRoleNotification(roleNotification *models.RoleNotification) (err error)
	CreatePostNotification(postNotification *models.PostNotification) (err error)
	DeleteAllRoleNotifications(userID int64) (err error)
	DeleteAllPostNotifications(userID int64) (err error)
	GetRoleNotifications(userID int64) (roleNotifications []models.RoleNotification, err error)
	GetPostNotifications(userID int64) (postNotifications []models.PostNotification, err error)
}
