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
