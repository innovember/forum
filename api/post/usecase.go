package post

import (
	"github.com/innovember/forum/api/models"
)

type PostUsecase interface {
	Create(post *models.Post, categories []string) (newPost *models.Post, status int, err error)
	GetAllPosts(userID int64) (posts []models.Post, status int, err error)
}
