package service

import (
	"errors"
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type GiveScoreService interface {
	QueryGiveScore(giveScoreThemeId int64) ([]*model.GiveScore, error)
	QueryGiveTeacherScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error)
	QueryGiveOneselfScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error)
	QueryGiveMutualScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error)
	AddGiveScore(giveScore model.GiveScore, files *multipart.Form) error
	EditGiveScore(giveScore model.GiveScore) error
	DelGiveSCore(giveScoreId int64) error
}

type giveScoreService struct {
	*Service
	giveScoreRepository repository.GiveScoreRepository
}

func NewGiveScoreService(service *Service, giveScoreRepository repository.GiveScoreRepository) GiveScoreService {
	return &giveScoreService{
		Service:             service,
		giveScoreRepository: giveScoreRepository,
	}
}

func (s *giveScoreService) QueryGiveScore(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	return s.giveScoreRepository.QueryGiveScore(giveScoreThemeId)
}

func (s *giveScoreService) QueryGiveTeacherScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	return s.giveScoreRepository.QueryGiveTeacherScoreByGiveScoreThemeId(giveScoreThemeId)
}

func (s *giveScoreService) QueryGiveOneselfScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	return s.giveScoreRepository.QueryGiveOneselfScoreByGiveScoreThemeId(giveScoreThemeId)
}

func (s *giveScoreService) QueryGiveMutualScoreByGiveScoreThemeId(giveScoreThemeId int64) ([]*model.GiveScore, error) {
	return s.giveScoreRepository.QueryGiveMutualScoreByGiveScoreThemeId(giveScoreThemeId)
}

func (s *giveScoreService) AddGiveScore(giveScore model.GiveScore, files *multipart.Form) error {
	return s.giveScoreRepository.AddGiveScore(giveScore, files)
}

func (s *giveScoreService) EditGiveScore(giveScore model.GiveScore) error {
	return s.giveScoreRepository.EditGiveScore(giveScore)
}

func (s *giveScoreService) DelGiveSCore(giveScoreId int64) error {
	if giveScoreId == 0 {
		return errors.New("参数错误")
	}
	return s.giveScoreRepository.DelGiveSCore(giveScoreId)
}
