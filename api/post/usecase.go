package post

import (
	"github.com/innovember/forum/api/models"
)

type PostUsecase interface {
	Create(post *models.Post, categories []string) (newPost *models.Post, status int, err error)
	GetAllPosts(userID int64) (posts []models.Post, status int, err error)
	GetPostByID(userID int64, postID int64) (post *models.Post, status int, err error)
	GetPostsByCategories(categories []string, userID int64) (posts []models.Post, status int, err error)
	GetPostsByRating(orderBy string, userID int64) (posts []models.Post, status int, err error)
	GetPostsByDate(orderBy string, userID int64) (posts []models.Post, status int, err error)
	GetAllPostsByAuthorID(authorID int64, userID int64) (posts []models.Post, status int, err error)
	GetRatedPostsByUser(userID int64, orderBy string, requestorID int64) (posts []models.Post, status int, err error)
	Update(post *models.Post) (editedPost *models.Post, status int, err error)
	Delete(postID int64) (status int, err error)
}

type CategoryUsecase interface {
	GetAllCategories() (categories []models.Category, status int, err error)
	Update(postID int64, categories []string) (err error)
	DeleteFromPostCategoriesBridge(postID int64) (err error)
}

type RateUsecase interface {
	RatePost(postID int64, userID int64, vote int) (int64, error)
	GetRating(postID int64, userID int64) (rating int, userRating int, err error)
	IsRatedBefore(postID int64, userID int64, vote int) (bool, error)
	DeleteRateFromPost(postID int64, userID int64, vote int) error
}

type CommentUsecase interface {
	Create(userID int64, comment *models.Comment) (newComment *models.Comment, status int, err error)
	GetCommentsByPostID(postID int64) (comments []models.Comment, status int, err error)
	GetCommentsByAuthorID(authorID int64) (comments []models.Comment, status int, err error)
	Update(comment *models.Comment) (editedComment *models.Comment, status int, err error)
	GetCommentByID(commentID int64) (comment *models.Comment, status int, err error)
	Delete(commentID int64) (status int, err error)
}

type NotificationUsecase interface {
	Create(notification *models.Notification) (newNotification *models.Notification, status int, err error)
	DeleteAllNotifications(receiverID int64) (status int, err error)
}
