package service

import (
	"errors"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type ChapterService interface {
	AddChapter(chapter model.Chapter) error
	EditChapter(chapter model.Chapter) error
	DelChapter(chapter model.Chapter) error
	QueryChapterByCourseId(CourseId int) ([]*model.ChapterView, error)
}

type chapterService struct {
	*Service
	chapterRepository repository.ChapterRepository
}

func NewChapterService(service *Service, chapterRepository repository.ChapterRepository) ChapterService {
	return &chapterService{
		Service:           service,
		chapterRepository: chapterRepository,
	}
}

func (s *chapterService) AddChapter(chapter model.Chapter) error {
	return s.chapterRepository.AddChapter(chapter)
}
func (s *chapterService) EditChapter(chapter model.Chapter) error {
	if chapter.ChapterId == 0 {
		return errors.New("参数错误")
	}
	return s.chapterRepository.EditChapter(chapter)

}
func (s *chapterService) DelChapter(chapter model.Chapter) error {
	if chapter.ChapterId == 0 {
		return errors.New("参数错误")
	}
	return s.chapterRepository.DelChapter(chapter)
}
func (s *chapterService) QueryChapterByCourseId(CourseId int) ([]*model.ChapterView, error) {
	if CourseId == 0 {
		return nil, errors.New("参数错误")
	}
	return s.chapterRepository.QueryChapterByCourseId(CourseId)
}
