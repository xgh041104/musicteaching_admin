package service

import (
	"errors"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type GiveScoreThemeService interface {
	QueryGiveScoreThemeByGiveScoreCategoryId(giveScoreCatgeoryId int64) ([]*model.GiveScoreTheme, error)
	QueryGiveScoreTheme(schoolId int64) ([]*model.GiveScoreThemeView, error)
	AddGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error
	EditGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error
	DelGiveScoreTheme(giveScoreThemeId int64) error
}

type giveScoreThemeService struct {
	*Service
	giveScoreThemeRepository repository.GiveScoreThemeRepository
}

func NewGiveScoreThemeService(service *Service, giveScoreThemeRepository repository.GiveScoreThemeRepository) GiveScoreThemeService {
	return &giveScoreThemeService{
		Service:                  service,
		giveScoreThemeRepository: giveScoreThemeRepository,
	}
}

func (s *giveScoreThemeService) QueryGiveScoreThemeByGiveScoreCategoryId(giveScoreCatgeoryId int64) ([]*model.GiveScoreTheme, error) {
	return s.giveScoreThemeRepository.QueryGiveScoreThemeByGiveScoreCategoryId(giveScoreCatgeoryId)
}

func (s *giveScoreThemeService) QueryGiveScoreTheme(schoolId int64) ([]*model.GiveScoreThemeView, error) {
	return s.giveScoreThemeRepository.QueryGiveScoreTheme(schoolId)
}

func (s *giveScoreThemeService) AddGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error {

	return s.giveScoreThemeRepository.AddGiveScoreTheme(giveScoreTheme)
}

func (s *giveScoreThemeService) EditGiveScoreTheme(giveScoreTheme model.GiveScoreTheme) error {
	if giveScoreTheme.GiveScoreCategoryId == 0 {
		return errors.New("参数错误")
	}

	return s.giveScoreThemeRepository.EditGiveScoreTheme(giveScoreTheme)
}

func (s *giveScoreThemeService) DelGiveScoreTheme(giveScoreThemeId int64) error {
	if giveScoreThemeId == 0 {
		return errors.New("参数错误")
	}

	return s.giveScoreThemeRepository.DelGiveScoreTheme(giveScoreThemeId)
}
