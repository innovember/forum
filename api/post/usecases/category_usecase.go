package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/post"
)

type CategoryUsecase struct {
	categoryRepo post.CategoryRepository
}

func NewCategoryUsecase(repo post.CategoryRepository) post.CategoryUsecase {
	return &CategoryUsecase{categoryRepo: repo}
}

func (cu *CategoryUsecase) GetAllCategories() (categories []models.Category, status int, err error) {
	if categories, status, err = cu.categoryRepo.GetAllCategories(); err != nil {
		return nil, status, err
	}
	return categories, status, nil
}

func (cu *CategoryUsecase) Update(postID int64, categories []string) (err error) {
	if err = cu.categoryRepo.Update(postID, categories); err != nil {
		return err
	}
	return nil
}

func (cu *CategoryUsecase) DeleteFromPostCategoriesBridge(postID int64) (err error) {
	if err = cu.categoryRepo.DeleteFromPostCategoriesBridge(postID); err != nil {
		return err
	}
	return nil
}

func (cu *CategoryUsecase) DeleteCategoryByID(categoryID int64) (err error) {
	if err = cu.categoryRepo.DeleteCategoryByID(categoryID); err != nil {
		return err
	}
	return nil
}

func (cu *CategoryUsecase) CreateNewCategory(category string) (err error) {
	if err = cu.categoryRepo.CreateNewCategory(category); err != nil {
		return err
	}
	return nil
}
