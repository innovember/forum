package post

import (
	"github.com/innovember/forum/api/models"
)

type PostRepository interface {
	Create(post *models.Post, categories []string) (newPost *models.Post, status int, err error)
	GetAllPosts(userID int64) (posts []models.Post, status int, err error)
	GetPostByID(userID int64, postID int64) (post *models.Post, status int, err error)
	GetCategories(post *models.Post) (status int, err error)
	GetAuthor(post *models.Post) (status int, err error)
	GetPostsByCategories(categories []string, userID int64) (posts []models.Post, status int, err error)
	GetPostsByRating(orderBy string, userID int64) (posts []models.Post, status int, err error)
	GetPostsByDate(orderBy string, userID int64) (posts []models.Post, status int, err error)
	GetAllPostsByAuthorID(authorID int64) (posts []models.Post, status int, err error)
	GetRatedPostsByUser(userID int64, orderBy string) (posts []models.Post, status int, err error)
}

type CategoryRepository interface {
	Create(postID int64, categories []string) (err error)
	GetAllCategories() (categories []models.Category, status int, err error)
}

type RateRepository interface {
	RatePost(postID int64, userID int64, vote int) error
	GetPostRating(postID int64, userID int64) (rating int, userRating int, err error)
	IsRatedBefore(postID int64, userID int64, vote int) (bool, error)
	DeleteRateFromPost(postID int64, userID int64, vote int) error
}

type CommentRepository interface {
	Create(userID int64, comment *models.Comment) (newComment *models.Comment, status int, err error)
	GetCommentsByPostID(postID int64) (comments []models.Comment, status int, err error)
	GetAuthor(comment *models.Comment) (status int, err error)
	GetCommentsByAuthorID(authorID int64) (comments []models.Comment, status int, err error)
	GetCommentsNumberByPostID(postID int64) (commentsNumber int, err error)
}
