package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
)

type PostUsecase struct {
	postRepo post.PostRepository
}

func NewPostUsecase(repo post.PostRepository) post.PostUsecase {
	return &PostUsecase{postRepo: repo}
}

func (pu *PostUsecase) Create(post *models.Post, categories []string) (newPost *models.Post, status int, err error) {
	if newPost, status, err = pu.postRepo.Create(post, categories); err != nil {
		return nil, status, err
	}
	return newPost, status, err
}

func (pu *PostUsecase) GetAllPosts(userID int64) (posts []models.Post, status int, err error) {
	if posts, status, err = pu.postRepo.GetAllPosts(userID); err != nil {
		return nil, status, err
	}
	return posts, status, nil
}

func (pu *PostUsecase) GetPostByID(userID int64, postID int64) (post *models.Post, status int, err error) {
	if post, status, err = pu.postRepo.GetPostByID(userID, postID); err != nil {
		return nil, status, err
	}
	return post, status, nil
}
