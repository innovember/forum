package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/user"
	"net/http"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{userRepo: repo}
}

func (uu *UserUsecase) Create(user *models.User) (status int, err error) {
	if status, err = uu.userRepo.CheckByUsernameOrEmail(user); err != nil {
		return status, err
	}
	if status, err = uu.userRepo.Create(user); err != nil {
		return status, err
	}
	return status, nil
}

func (uu *UserUsecase) GetAllUsers() (users []models.User, status int, err error) {
	if users, err = uu.userRepo.GetAllUsers(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return users, http.StatusOK, nil
}
