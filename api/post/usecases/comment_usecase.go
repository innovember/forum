package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
)

type CommentUsecase struct {
	commentRepo post.CommentRepository
}

func NewCommentUsecase(repo post.CommentRepository) post.CommentUsecase {
	return &CommentUsecase{commentRepo: repo}
}
func (cu *CommentUsecase) Create(userID int64, comment *models.Comment) (newComment *models.Comment, status int, err error) {
	if newComment, status, err = cu.commentRepo.Create(userID, comment); err != nil {
		return nil, status, err
	}
	return newComment, status, err
}
func (cu *CommentUsecase) GetCommentsByPostID(postID int64) (comments []models.Comment, status int, err error) {
	if comments, status, err = cu.commentRepo.GetCommentsByPostID(postID); err != nil {
		return nil, status, err
	}
	return comments, status, err
}

func (cu *CommentUsecase) GetCommentsByAuthorID(authorID int64) (comments []models.Comment, status int, err error) {
	if comments, status, err = cu.commentRepo.GetCommentsByAuthorID(authorID); err != nil {
		return nil, status, err
	}
	return comments, status, err
}

func (cu *CommentUsecase) Update(comment *models.Comment) (editedComment *models.Comment, status int, err error) {
	if editedComment, status, err = cu.commentRepo.Update(comment); err != nil {
		return nil, status, err
	}
	return editedComment, status, err
}
