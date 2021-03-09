package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/user"
)

type ModeratorUsecase struct {
	moderatorRepo user.ModeratorRepository
}

func NewModeratorUsecase(repo user.ModeratorRepository) user.ModeratorUsecase {
	return &ModeratorUsecase{moderatorRepo: repo}
}

func (mu *ModeratorUsecase) CreatePostReport(postReport *models.PostReport) (err error) {
	if err = mu.moderatorRepo.CreatePostReport(postReport); err != nil {
		return err
	}
	return nil
}
func (mu *ModeratorUsecase) DeletePostReport(postReportID int64) (err error) {
	if err = mu.moderatorRepo.DeletePostReport(postReportID); err != nil {
		return err
	}
	return nil
}

func (mu *ModeratorUsecase) GetMyReports(moderatorID int64) (postReports []models.PostReport, err error) {
	if postReports, err = mu.moderatorRepo.GetMyReports(moderatorID); err != nil {
		return nil, err
	}
	return postReports, nil
}

func (mu *ModeratorUsecase) ApprovePost(postID int64) (err error) {
	if err = mu.moderatorRepo.ApprovePost(postID); err != nil {
		return err
	}
	return nil
}

func (mu *ModeratorUsecase) GetAllUnapprovedPosts() (posts []models.Post, err error) {
	if posts, err = mu.moderatorRepo.GetAllUnapprovedPosts(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (mu *ModeratorUsecase) BanPost(postID int64, bans []string) (err error) {
	if err = mu.moderatorRepo.BanPost(postID, bans); err != nil {
		return err
	}
	return nil
}

func (mu *ModeratorUsecase) GetPostReportByID(postReportID int64) (postReport *models.PostReport, err error) {
	if postReport, err = mu.moderatorRepo.GetPostReportByID(postReportID); err != nil {
		return nil, err
	}
	return postReport, nil
}
