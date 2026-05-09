package repository

import (
	"ai_summary_project/internal/model"
	"context"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(ctx context.Context, book *model.Book) error
	GetByID(ctx context.Context, id uint) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*model.Book, error)
	AddCourseCount(ctx context.Context, bookID uint, delta int) error
}

func NewBookRepository(
	repository *Repository,
) BookRepository {
	return &bookRepo{
		Repository: repository,
	}
}

type bookRepo struct {
	*Repository
}

func (r *bookRepo) Create(ctx context.Context, book *model.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepo) GetByID(ctx context.Context, id uint) (*model.Book, error) {
	var book model.Book
	if err := r.db.WithContext(ctx).First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
func (r *bookRepo) AddCourseCount(ctx context.Context, bookID uint, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Book{}).
		Where("id = ?", bookID).
		UpdateColumn("course_count", gorm.Expr("course_count + ?", delta)).Error
}

func (r *bookRepo) Update(ctx context.Context, book *model.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *bookRepo) Delete(ctx context.Context, id uint) error {
	// 软删除
	return r.db.WithContext(ctx).Delete(&model.Book{}, id).Error
}

func (r *bookRepo) List(ctx context.Context) ([]*model.Book, error) {
	var books []*model.Book
	if err := r.db.WithContext(ctx).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
