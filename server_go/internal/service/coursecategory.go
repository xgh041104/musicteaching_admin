package service

import (
	"errors"
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type CourseCategoryService interface {
	AddCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error
	EditCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error
	DelCourseCategory(courseCategory model.CourseCategory) error
	QueryCourseCategoryTree(schoolId int) ([]*model.CourseCategoryTreeList, error)
}

type courseCategoryService struct {
	*Service
	courseCategoryRepository repository.CourseCategoryRepository
}

func NewCourseCategoryService(service *Service, courseCategoryRepository repository.CourseCategoryRepository) CourseCategoryService {
	return &courseCategoryService{
		Service:                  service,
		courseCategoryRepository: courseCategoryRepository,
	}
}

func (s *courseCategoryService) AddCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error {

	return s.courseCategoryRepository.AddCourseCategory(courseCategory, files)
}

func (s *courseCategoryService) EditCourseCategory(courseCategory model.CourseCategory, files *multipart.Form) error {
	if courseCategory.CourseCategoryId == 0 {
		return errors.New("参数错误")
	}
	return s.courseCategoryRepository.EditCourseCategory(courseCategory, files)
}
func (s *courseCategoryService) DelCourseCategory(courseCategory model.CourseCategory) error {
	if courseCategory.CourseCategoryId == 0 {
		return errors.New("参数错误")
	}
	return s.courseCategoryRepository.DelCourseCategory(courseCategory)
}

func (s *courseCategoryService) QueryCourseCategoryTree(schoolId int) ([]*model.CourseCategoryTreeList, error) {
	return s.courseCategoryRepository.QueryCourseCategoryTree(schoolId)
}
