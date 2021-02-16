package usecases

import (
	"github.com/innovember/forum/api/post"
)

type RateUsecase struct {
	rateRepo post.RateRepository
}

func NewRateUsecase(repo post.RateRepository) post.RateUsecase {
	return &RateUsecase{rateRepo: repo}
}

func (ru *RateUsecase) RatePost(postID int64, userID int64, vote int) (rateID int64, err error) {
	if rateID, err = ru.rateRepo.RatePost(postID, userID, vote); err != nil {
		return 0, err
	}
	return rateID, nil
}
func (ru *RateUsecase) GetRating(postID int64, userID int64) (rating int, userRating int, err error) {
	if rating, userRating, err = ru.rateRepo.GetPostRating(postID, userID); err != nil {
		return 0, 0, err
	}
	return rating, userRating, nil
}
func (ru *RateUsecase) IsRatedBefore(postID int64, userID int64, vote int) (bool, error) {
	var (
		isRated bool
		err     error
	)
	if isRated, err = ru.rateRepo.IsRatedBefore(postID, userID, vote); err != nil {
		return false, err
	}
	if isRated {
		return true, nil
	}
	return false, nil
}

func (ru *RateUsecase) DeleteRateFromPost(postID int64, userID int64, vote int) error {
	if err := ru.rateRepo.DeleteRateFromPost(postID, userID, vote); err != nil {
		return err
	}
	return nil
}

func (ru *RateUsecase) DeleteRatesByPostID(postID int64) (err error) {
	if err := ru.rateRepo.DeleteRatesByPostID(postID); err != nil {
		return err
	}
	return nil
}
