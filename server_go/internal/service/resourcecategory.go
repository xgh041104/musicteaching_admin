package service

import (
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type ResourceCategoryService interface {
	AddResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error
	EditResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error
	DelResourceCategory(resourceCategory model.ResourcesCategory) error
	QueryResourceCategoryTree(schoolId int) ([]*model.ResourcesCategoryTreeList, error)
	QueryResourceCategoryParentNodeByParentId(resourceCategoryParentId int) (*model.ResourcesCategoryView, error)
	QueryResourceCategoryChildNodesById(resourceCategoryId int, schoolId int) ([]*model.ResourcesCategoryView, error)
}

type resourceCategoryService struct {
	*Service
	resourceCategoryRepository repository.ResourceCategoryRepository
}

func NewResourceCategoryService(service *Service, resourceCategoryRepository repository.ResourceCategoryRepository) ResourceCategoryService {
	return &resourceCategoryService{
		Service:                    service,
		resourceCategoryRepository: resourceCategoryRepository,
	}
}

func (s *resourceCategoryService) AddResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error {

	return s.resourceCategoryRepository.AddResourceCategory(resourceCategory, files)
}

func (s *resourceCategoryService) EditResourceCategory(resourceCategory model.ResourcesCategory, files *multipart.Form) error {

	return s.resourceCategoryRepository.EditResourceCategory(resourceCategory, files)
}
func (s *resourceCategoryService) DelResourceCategory(resourceCategory model.ResourcesCategory) error {

	return s.resourceCategoryRepository.DelResourceCategory(resourceCategory)
}
func (s *resourceCategoryService) QueryResourceCategoryTree(schoolId int) ([]*model.ResourcesCategoryTreeList, error) {

	return s.resourceCategoryRepository.QueryResourceCategoryTree(schoolId)
}

func (s *resourceCategoryService) QueryResourceCategoryParentNodeByParentId(resourceCategoryParentId int) (*model.ResourcesCategoryView, error) {
	return s.resourceCategoryRepository.QueryResourceCategoryParentNodeByParentId(resourceCategoryParentId)
}

func (s *resourceCategoryService) QueryResourceCategoryChildNodesById(resourceCategoryId int, schoolId int) ([]*model.ResourcesCategoryView, error) {
	return s.resourceCategoryRepository.QueryResourceCategoryChildNodesById(resourceCategoryId, schoolId)
}
