package service

import (
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type SectionService interface {
	EditerUploadFile(files *multipart.Form) (string, error)
	AddSection(section model.SectionEdit, files *multipart.Form) error
	DelSection(section model.Section) error
	EditSection(section model.SectionEdit, files *multipart.Form) error
	QuerySectionBySectionId(SectionId int) (model.SectionView, error)
	QuerySectionByChapterId(chapterId int) ([]*model.SectionModel, error)
}

type sectionService struct {
	*Service
	sectionRepository repository.SectionRepository
}

func NewSectionService(service *Service, sectionRepository repository.SectionRepository) SectionService {
	return &sectionService{
		Service:           service,
		sectionRepository: sectionRepository,
	}
}

func (s *sectionService) EditerUploadFile(files *multipart.Form) (string, error) {

	return s.sectionRepository.EditerUploadFile(files)
}
func (s *sectionService) AddSection(section model.SectionEdit, files *multipart.Form) error {
	return s.sectionRepository.AddSection(section, files)
}

func (s *sectionService) DelSection(section model.Section) error {
	return s.sectionRepository.DelSection(section)
}
func (s *sectionService) EditSection(section model.SectionEdit, files *multipart.Form) error {

	return s.sectionRepository.EditSection(section, files)
}

func (s *sectionService) QuerySectionBySectionId(SectionId int) (model.SectionView, error) {
	return s.sectionRepository.QuerySectionBySectionId(SectionId)
}

func (s *sectionService) QuerySectionByChapterId(chapterId int) ([]*model.SectionModel, error) {
	return s.sectionRepository.QuerySectionByChapterId(chapterId)
}
