package service

import (
	"errors"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type AdminService interface {
	GetSchoolAdmin() ([]*model.Admin, error)
	AddSchoolAdmin(admin model.Admin) error
	DelSchoolAdmin(adminId int64) error
	UpdateSchoolAdmin(admin model.Admin) error
}

type adminService struct {
	*Service
	adminRepository repository.AdminRepository
}

func NewAdminService(service *Service, adminRepository repository.AdminRepository) AdminService {
	return &adminService{
		Service:         service,
		adminRepository: adminRepository,
	}
}

func (s *adminService) AddSchoolAdmin(admin model.Admin) error {
	return s.adminRepository.AddSchoolAdmin(admin)
}

func (s *adminService) GetSchoolAdmin() ([]*model.Admin, error) {

	return s.adminRepository.GetSchoolAdmin()
}

func (s *adminService) DelSchoolAdmin(adminId int64) error {

	if adminId == 0 {
		return errors.New("参数错误")
	}
	return s.adminRepository.DelSchoolAdmin(adminId)
}

func (s *adminService) UpdateSchoolAdmin(admin model.Admin) error {
	if admin.AdminId == 0 {
		return errors.New("参数错误")
	}
	return s.adminRepository.UpdateSchoolAdmin(admin)
}
