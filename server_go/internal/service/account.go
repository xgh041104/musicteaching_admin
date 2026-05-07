package service

import (
	"errors"
	"mime/multipart"
	"server_go/internal/model"
	"server_go/internal/repository"
)

type AccountService interface {
	AddAccount(account model.Account, files *multipart.Form) error
	EditAccount(account model.Account, files *multipart.Form) error
	DelAccount(account model.Account) error
	QueryAllAccount() ([]*model.AccountView, error)
}

type accountService struct {
	*Service
	accountRepository repository.AccountRepository
}

func NewAccountService(service *Service, accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		Service:           service,
		accountRepository: accountRepository,
	}
}

func (s *accountService) AddAccount(account model.Account, files *multipart.Form) error {

	return s.accountRepository.AddAccount(account, files)
}

func (s *accountService) EditAccount(account model.Account, files *multipart.Form) error {

	if account.AccountId == 0 {
		return errors.New("参数错误")
	}
	return s.accountRepository.EditAccount(account, files)
}

func (s *accountService) DelAccount(account model.Account) error {
	if account.AccountId == 0 {
		return errors.New("参数错误")
	}
	return s.accountRepository.DelAccount(account)
}

func (s *accountService) QueryAllAccount() ([]*model.AccountView, error) {

	return s.accountRepository.QueryAllAccount()
}
