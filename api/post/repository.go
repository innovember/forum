package post

import (
	"github.com/innovember/forum/api/models"
)

type PostRepository interface {
	Create(post *models.Post, categories []string) (newPost *models.Post, status int, err error)
	GetAllPosts(userID int64) (posts []models.Post, status int, err error)
	GetPostByID(userID int64, postID int64) (post *models.Post, status int, err error)
	// GetPostsByCategories()
	// GetPostsByRating()
	// GetPostsByDate()
}

type CategoryRepository interface {
	Create(postID int64, categories []string) (err error)
}

type RateRepository interface {
	RatePost(postID int64, userID int64, vote int) error
	GetPostRating(postID int64, userID int64) (rating int, userRating int, err error)
}
