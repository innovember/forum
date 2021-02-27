package usecases

import (
	"github.com/innovember/forum/api/models"
	"github.com/innovember/forum/api/user"
)

type AdminUsecase struct {
	adminRepo user.AdminRepository
}

func NewAdminUsecase(repo user.AdminRepository) user.AdminUsecase {
	return &AdminUsecase{adminRepo: repo}
}

func (au *AdminUsecase) UpgradeRole(userID int64) (err error) {
	if err = au.adminRepo.UpgradeRole(userID); err != nil {
		return err
	}
	return nil
}

func (au *AdminUsecase) GetAllRoleRequests() (roleRequests []models.RoleRequest, err error) {
	if roleRequests, err = au.adminRepo.GetAllRoleRequests(); err != nil {
		return nil, err
	}
	return roleRequests, nil
}

func (au *AdminUsecase) DeleteRoleRequest(userID int64) (err error) {
	if err = au.adminRepo.DeleteRoleRequest(userID); err != nil {
		return err
	}
	return nil
}
