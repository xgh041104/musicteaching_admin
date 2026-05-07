package service

import (
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type ResourceService interface {
	AddResource(resource model.Resources, files *multipart.Form) error
	DelResource(resource model.Resources) error
	QueryResourceByResourceCategoryId(schoolId int, lecturerCommonUserId int, ResourceCategoryId int) ([]*model.ResourcesView, error)
	QueryResourceById(resourceId int) (*model.ResourcesView, error)
}

type resourceService struct {
	*Service
	resourceRepository repository.ResourceRepository
}

func NewResourceService(service *Service, resourceRepository repository.ResourceRepository) ResourceService {
	return &resourceService{
		Service:            service,
		resourceRepository: resourceRepository,
	}

}

func (s *resourceService) AddResource(resource model.Resources, files *multipart.Form) error {

	return s.resourceRepository.AddResource(resource, files)
}

func (s *resourceService) DelResource(resource model.Resources) error {

	return s.resourceRepository.DelResource(resource)
}

func (s *resourceService) QueryResourceByResourceCategoryId(schoolId int, lecturerCommonUserId int, ResourceCategoryId int) ([]*model.ResourcesView, error) {

	return s.resourceRepository.QueryResourceByResourceCategoryId(schoolId, lecturerCommonUserId, ResourceCategoryId)
}

func (s *resourceService) QueryResourceById(resourceId int) (*model.ResourcesView, error) {

	return s.resourceRepository.QueryResourceById(resourceId)
}
