package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
)

type BanUsecase struct {
	banRepo post.BanRepository
}

func NewBanUsecase(repo post.BanRepository) post.BanUsecase {
	return &BanUsecase{banRepo: repo}
}

func (bu *BanUsecase) GetAllCategories() (categories []models.Category, status int, err error) {
	if categories, status, err = bu.banRepo.GetAllCategories(); err != nil {
		return nil, status, err
	}
	return categories, status, nil
}

func (bu *BanUsecase) Update(postID int64, categories []string) (err error) {
	if err = bu.banRepo.Update(postID, categories); err != nil {
		return err
	}
	return nil
}

func (bu *BanUsecase) DeleteFromPostCategoriesBridge(postID int64) (err error) {
	if err = bu.banRepo.DeleteFromPostCategoriesBridge(postID); err != nil {
		return err
	}
	return nil
}

func (bu *BanUsecase) DeleteCategoryByID(categoryID int64) (err error) {
	if err = bu.banRepo.DeleteCategoryByID(categoryID); err != nil {
		return err
	}
	return nil
}

func (bu *BanUsecase) CreateNewCategory(category string) (err error) {
	if err = bu.banRepo.CreateNewCategory(category); err != nil {
		return err
	}
	return nil
}
