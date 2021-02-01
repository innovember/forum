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

func (ru *RateUsecase) RatePost(postID int64, userID int64, vote int) error {
	if err := ru.rateRepo.RatePost(postID, userID, vote); err != nil {
		return err
	}
	return nil
}
func (ru *RateUsecase) GetRating(postID int64, userID int64) (rating int, userRating int, err error) {
	if rating, userRating, err = ru.rateRepo.GetPostRating(postID, userID); err != nil {
		return 0, 0, err
	}
	return rating, userRating, nil
}
