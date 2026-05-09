package service

import (
	v1 "ai_summary_project/api/v1"
	"ai_summary_project/internal/model"
	"ai_summary_project/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersService interface {
	GetUsers(ctx context.Context, id int64) (*model.Users, error)
	Login(ctx context.Context, req *v1.Login) (*v1.LoginResponse, error)
}

func NewUsersService(
	service *Service,
	usersRepository repository.UsersRepository,
) UsersService {
	return &usersService{
		Service:         service,
		usersRepository: usersRepository,
	}
}

type usersService struct {
	*Service
	usersRepository repository.UsersRepository
}

func (s *usersService) GetUsers(ctx context.Context, id int64) (*model.Users, error) {
	return s.usersRepository.GetUsers(ctx, id)
}
func (s *usersService) Login(ctx context.Context, req *v1.Login) (*v1.LoginResponse, error) {
	if req.UserAccount == "" || req.UserPwd == "" {
		return nil, v1.ErrBadRequest
	}
	user, err := s.usersRepository.GetUsersByAccount(ctx, req.UserAccount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, v1.ErrInternalServerError
	}
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println(user)
	if user.UserType != model.UserType {
		return nil, errors.New("用户类型错误")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.UserPwd), []byte(req.UserPwd))
	if err != nil {
		return nil, errors.New("加密失败")
	}
	var users v1.User
	users.ID = user.ID
	users.UserAccount = user.UserAccount
	users.UserTrueName = user.UserTrueName
	users.UserType = int(user.UserType)
	users.UserPwd = user.UserPwd
	users.SchoolId = user.SchoolId
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println(users)
	token, _ := s.jwt.GenToken(user.ID, time.Now().Add(time.Hour*24*7), int(user.UserType))
	var info = v1.LoginResponse{
		Data:  &users,
		Token: token,
	}
	return &info, nil
}
