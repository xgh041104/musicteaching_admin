package service

import (
	"errors"
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type SchoolService interface {
	AddSchool(school model.School, files *multipart.Form) error
	EditSchool(school model.School, files *multipart.Form) error
	DelSchool(school model.School) error
	QuerySchoolAll() ([]*model.SchoolView, error)
}

type schoolService struct {
	*Service
	schoolRepository repository.SchoolRepository
}

func NewSchoolService(service *Service, schoolRepository repository.SchoolRepository) SchoolService {
	return &schoolService{
		Service:          service,
		schoolRepository: schoolRepository,
	}
}

func (s *schoolService) AddSchool(school model.School, files *multipart.Form) error {

	return s.schoolRepository.AddSchool(school, files)
}

func (s *schoolService) EditSchool(school model.School, files *multipart.Form) error {

	if school.SchoolId == 0 {
		return errors.New("参数错误")
	}
	return s.schoolRepository.EditSchool(school, files)
}

func (s *schoolService) DelSchool(school model.School) error {
	if school.SchoolId == 0 {
		return errors.New("参数错误")
	}
	return s.schoolRepository.DelSchool(school)
}

func (s *schoolService) QuerySchoolAll() ([]*model.SchoolView, error) {

	return s.schoolRepository.QuerySchoolAll()
}
