package usecases

import (
	"github.com/innovember/forum/api/post"
)

type RateCommentUsecase struct {
	rateCommentRepo post.RateCommentRepository
}

func NewRateCommentUsecase(repo post.RateCommentRepository) post.RateCommentUsecase {
	return &RateCommentUsecase{rateCommentRepo: repo}
}

func (ru *RateCommentUsecase) RateComment(commentID int64, userID int64, vote int, postID int64) (rateID int64, err error) {
	if rateID, err = ru.rateCommentRepo.RateComment(commentID, userID, vote, postID); err != nil {
		return 0, err
	}
	return rateID, nil
}
func (ru *RateCommentUsecase) GetCommentRating(commentID int64, userID int64) (rating int, userRating int, err error) {
	if rating, userRating, err = ru.rateCommentRepo.GetCommentRating(commentID, userID); err != nil {
		return 0, 0, err
	}
	return rating, userRating, nil
}
func (ru *RateCommentUsecase) IsRatedBefore(commentID int64, userID int64, vote int) (bool, error) {
	var (
		isRated bool
		err     error
	)
	if isRated, err = ru.rateCommentRepo.IsRatedBefore(commentID, userID, vote); err != nil {
		return false, err
	}
	if isRated {
		return true, nil
	}
	return false, nil
}

func (ru *RateCommentUsecase) DeleteRateFromComment(commentID int64, userID int64, vote int) error {
	if err := ru.rateCommentRepo.DeleteRateFromComment(commentID, userID, vote); err != nil {
		return err
	}
	return nil
}

func (ru *RateCommentUsecase) DeleteRatesByCommentID(commentID int64) (err error) {
	if err := ru.rateCommentRepo.DeleteRatesByCommentID(commentID); err != nil {
		return err
	}
	return nil
}
