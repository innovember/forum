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
}

type CategoryUsecase interface {
	GetAllCategories() (categories []models.Category, status int, err error)
}

type RateUsecase interface {
	RatePost(postID int64, userID int64, vote int) error
	GetRating(postID int64, userID int64) (rating int, userRating int, err error)
}
