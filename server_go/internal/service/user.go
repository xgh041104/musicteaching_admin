package service

import (
	"errors"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type UserService interface {
	LoginUser(userparam model.LoginReq) (*model.LoginUser, error)
	EditUserPwd(userparam model.LoginUser) error
	GetUserBySchoolId(schoolId int64) ([]*model.CommonUser, error)
	AddCommonUser(commonUser model.CommonUser) error
	UpdateCommonUser(commonUser model.CommonUser) error
	DelCommonUser(commonUserId int64) error
}

type userService struct {
	*Service
	userRepository repository.UserRepository
}

func NewUserService(service *Service, userRepository repository.UserRepository) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
	}
}
func (s *userService) LoginUser(userparam model.LoginReq) (*model.LoginUser, error) {

	return s.userRepository.LoginUser(userparam)
}

func (s *userService) EditUserPwd(userparam model.LoginUser) error {
	if userparam.UserId == 0 {
		return errors.New("参数错误")
	}
	return s.userRepository.EditUserPwd(userparam)
}

func (s *userService) GetUserBySchoolId(schoolId int64) ([]*model.CommonUser, error) {

	if schoolId == 0 {
		return nil, errors.New("参数错误")
	}
	return s.userRepository.GetUserBySchoolId(schoolId)
}

func (s *userService) AddCommonUser(commonUser model.CommonUser) error {
	return s.userRepository.AddCommonUser(commonUser)
}

func (s *userService) UpdateCommonUser(commonUser model.CommonUser) error {
	if commonUser.CommonUserId == 0 {
		return errors.New("参数错误")
	}
	return s.userRepository.UpdateCommonUser(commonUser)
}

func (s *userService) DelCommonUser(commonUserId int64) error {
	if commonUserId == 0 {
		return errors.New("参数错误")
	}
	return s.userRepository.DelCommonUser(commonUserId)
}
