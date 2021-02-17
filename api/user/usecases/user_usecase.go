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

func (uu *UserUsecase) GetPassword(username string) (password string, status int, err error) {
	if password, status, err = uu.userRepo.GetPassword(username); err != nil {
		return "", status, err
	}
	return password, status, nil
}

func (uu *UserUsecase) FindUserByUsername(username string) (user *models.User, status int, err error) {
	if user, status, err = uu.userRepo.FindUserByUsername(username); err != nil {
		return nil, status, err
	}
	return user, status, nil
}

func (uu *UserUsecase) UpdateSession(userID int64, sessionValue string) (err error) {
	if err = uu.userRepo.UpdateSession(userID, sessionValue); err != nil {
		return err
	}
	return nil
}
func (uu *UserUsecase) ValidateSession(sessionValue string) (user *models.User, status int, err error) {
	if user, status, err = uu.userRepo.ValidateSession(sessionValue); err != nil {
		return nil, status, err
	}
	return user, status, nil
}

func (uu *UserUsecase) CheckSessionByUsername(username string) (status int, err error) {
	if status, err = uu.userRepo.CheckSessionByUsername(username); err != nil {
		return status, err
	}
	return status, nil
}
func (uu *UserUsecase) GetUserByID(userID int64) (user *models.User, err error) {
	if user, err = uu.userRepo.GetUserByID(userID); err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *UserUsecase) UpdateActivity(userID int64) (err error) {
	if err = uu.userRepo.UpdateActivity(userID); err != nil {
		return err
	}
	return nil
}
