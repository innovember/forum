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

func (au *AdminUsecase) UpgradeRole(requestID int64) (err error) {
	if err = au.adminRepo.UpgradeRole(requestID); err != nil {
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

func (au *AdminUsecase) DeleteRoleRequest(requestID int64) (err error) {
	if err = au.adminRepo.DeleteRoleRequest(requestID); err != nil {
		return err
	}
	return nil
}

func (au *AdminUsecase) GetAllPostReports() (postReports []models.PostReport, err error) {
	if postReports, err = au.adminRepo.GetAllPostReports(); err != nil {
		return nil, err
	}
	return postReports, nil
}

func (au *AdminUsecase) AcceptPostReport(postReportID int64) (err error) {
	if err = au.adminRepo.AcceptPostReport(postReportID); err != nil {
		return err
	}
	return nil
}

func (au *AdminUsecase) DismissPostReport(postReportID int64) (err error) {
	if err = au.adminRepo.DismissPostReport(postReportID); err != nil {
		return err
	}
	return nil
}

func (au *AdminUsecase) GetAllModerators() (moderators []models.User, err error) {
	if moderators, err = au.adminRepo.GetAllModerators(); err != nil {
		return nil, err
	}
	return moderators, nil
}

func (au *AdminUsecase) DemoteModerator(moderatorID int64) (err error) {
	if err = au.adminRepo.DemoteModerator(moderatorID); err != nil {
		return err
	}
	return nil
}
