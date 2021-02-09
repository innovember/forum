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

func (pu *PostUsecase) GetPostsByCategories(categories []string, userID int64) (posts []models.Post, status int, err error) {
	if posts, status, err = pu.postRepo.GetPostsByCategories(categories, userID); err != nil {
		return nil, status, err
	}
	return posts, status, nil
}

func (pu *PostUsecase) GetPostsByRating(orderBy string, userID int64) (posts []models.Post, status int, err error) {
	if posts, status, err = pu.postRepo.GetPostsByRating(orderBy, userID); err != nil {
		return nil, status, err
	}
	return posts, status, nil
}

func (pu *PostUsecase) GetPostsByDate(orderBy string, userID int64) (posts []models.Post, status int, err error) {
	if posts, status, err = pu.postRepo.GetPostsByDate(orderBy, userID); err != nil {
		return nil, status, err
	}
	return posts, status, nil
}

func (pu *PostUsecase) GetAllPostsByAuthorID(authorID int64) (posts []models.Post, status int, err error) {
	if posts, status, err = pu.postRepo.GetAllPostsByAuthorID(authorID); err != nil {
		return nil, status, err
	}
	return posts, status, nil
}
func (pu *PostUsecase) GetRatedPostsByUser(userID int64, orderBy string) (posts []models.Post, status int, err error) {
	if posts, status, err = pu.postRepo.GetRatedPostsByUser(userID, orderBy); err != nil {
		return nil, status, err
	}
	return posts, status, nil
}

func (pu *PostUsecase) Update(post *models.Post) (editedPost *models.Post, status int, err error) {
	if editedPost, status, err = pu.postRepo.Update(post); err != nil {
		return nil, status, err
	}
	return editedPost, status, err
}
func (pu *PostUsecase) Delete(postID int64) (status int, err error) {
	if status, err = pu.postRepo.Delete(postID); err != nil {
		return status, err
	}
	return status, nil
}
