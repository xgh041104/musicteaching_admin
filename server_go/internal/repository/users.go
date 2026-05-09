package repository

import (
	"ai_summary_project/internal/model"
	"context"
)

type UsersRepository interface {
	GetUsers(ctx context.Context, id int64) (*model.Users, error)
	GetUsersByAccount(ctx context.Context, account string) (*model.Users, error)
}

func NewUsersRepository(
	repository *Repository,
) UsersRepository {
	return &usersRepository{
		Repository: repository,
	}
}

type usersRepository struct {
	*Repository
}

func (r *usersRepository) GetUsers(ctx context.Context, id int64) (*model.Users, error) {
	var users model.Users

	return &users, nil
}
func (r *usersRepository) GetUsersByAccount(ctx context.Context, account string) (*model.Users, error) {
	var users model.Users
	err := r.db.Where("user_account = ?", account).First(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}
