package service

import (
	"errors"
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type GiveScoreCategoryService interface {
	QueryGiveScoreCategory(schoolId int64) ([]*model.GiveScoreCategoryView, error)
	QueryGiveScoreCategoryTree(schoolId int64) ([]*model.GiveScoreCategoryTreeList, error)
	QueryGiveScoreCategoryParentNodeByParentId(giveScoreCategoryParentId int64) (*model.GiveScoreCategoryView, error)
	QueryGiveScoreCategoryChildNodesById(giveScoreCategoryId int64, schoolId int64) ([]*model.GiveScoreCategoryView, error)
	AddGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error
	EditGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error
	DelGiveScoreCategory(giveScoreCategoryId int64) error
}

type giveScoreCategoryService struct {
	*Service
	giveScoreCategoryRepository repository.GiveScoreCategoryRepository
}

func NewGiveScoreCategoryService(service *Service, giveScoreCategoryRepository repository.GiveScoreCategoryRepository) GiveScoreCategoryService {
	return &giveScoreCategoryService{
		Service:                     service,
		giveScoreCategoryRepository: giveScoreCategoryRepository,
	}
}

func (s *giveScoreCategoryService) AddGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error {
	return s.giveScoreCategoryRepository.AddGiveScoreCategory(giveScoreCategory, files)
}

func (s *giveScoreCategoryService) DelGiveScoreCategory(giveScoreCategoryId int64) error {
	if giveScoreCategoryId == 0 {
		return errors.New("参数错误")
	}
	return s.giveScoreCategoryRepository.DelGiveScoreCategory(giveScoreCategoryId)
}

func (s *giveScoreCategoryService) EditGiveScoreCategory(giveScoreCategory model.GiveScoreCategory, files *multipart.Form) error {
	if giveScoreCategory.GiveScoreCategoryId == 0 {
		return errors.New("参数错误")
	}
	return s.giveScoreCategoryRepository.EditGiveScoreCategory(giveScoreCategory, files)
}

func (s *giveScoreCategoryService) QueryGiveScoreCategory(schoolId int64) ([]*model.GiveScoreCategoryView, error) {
	return s.giveScoreCategoryRepository.QueryGiveScoreCategory(schoolId)
}

func (s *giveScoreCategoryService) QueryGiveScoreCategoryTree(schoolId int64) ([]*model.GiveScoreCategoryTreeList, error) {
	return s.giveScoreCategoryRepository.QueryGiveScoreCategoryTree(schoolId)
}

func (s *giveScoreCategoryService) QueryGiveScoreCategoryParentNodeByParentId(giveScoreCategoryParentId int64) (*model.GiveScoreCategoryView, error) {
	return s.giveScoreCategoryRepository.QueryGiveScoreCategoryParentNodeByParentId(giveScoreCategoryParentId)
}
func (s *giveScoreCategoryService) QueryGiveScoreCategoryChildNodesById(giveScoreCategoryId int64, schoolId int64) ([]*model.GiveScoreCategoryView, error) {
	return s.giveScoreCategoryRepository.QueryGiveScoreCategoryChildNodesById(giveScoreCategoryId, schoolId)
}
