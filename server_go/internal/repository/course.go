package repository

import (
	"ai_summary_project/internal/model"
	"context"
)

type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
	GetByID(ctx context.Context, id uint) (*model.Course, error)
	Update(ctx context.Context, course *model.Course) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*model.Course, error)
	ListByBookID(ctx context.Context, bookID uint) ([]*model.Course, error)
}

func NewCourseRepository(
	repository *Repository,
) CourseRepository {
	return &courseRepo{
		Repository: repository,
	}
}

type courseRepo struct {
	*Repository
}

func (r *courseRepo) Create(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *courseRepo) GetByID(ctx context.Context, id uint) (*model.Course, error) {
	var course model.Course
	if err := r.db.WithContext(ctx).First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepo) Update(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *courseRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Course{}, id).Error
}

func (r *courseRepo) List(ctx context.Context) ([]*model.Course, error) {
	var courses []*model.Course
	if err := r.db.WithContext(ctx).Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepo) ListByBookID(ctx context.Context, bookID uint) ([]*model.Course, error) {
	var list []*model.Course
	if err := r.db.WithContext(ctx).Where("book_id = ?", bookID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
