package service

import (
	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/Vladislave00/CashFlow/app/pkg/repository"
)

type AccountService struct {
	repo           repository.Account
	accountingRepo repository.Accounting
}

func NewAccountService(repo repository.Account, accountingRepo repository.Accounting) *AccountService {
	return &AccountService{repo: repo, accountingRepo: accountingRepo}
}

func (s *AccountService) CreateAccount(userId, accountingId int, account models.CreateAccountInput) (int, error) {
	_, err := s.accountingRepo.GetById(userId, accountingId)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateAccount(accountingId, account)
}

func (s *AccountService) GetAll(userId, accountingId int) ([]models.Account, error) {
	return s.repo.GetAll(userId, accountingId)
}

func (s *AccountService) GetById(userId, accountId int) (models.Account, error) {
	return s.repo.GetById(userId, accountId)
}

func (s *AccountService) Delete(userId, accountId int) error {
	return s.repo.Delete(userId, accountId)
}

func (s *AccountService) Update(userId, accountId int, input models.UpdateAccountInput) error {
	return s.repo.Update(userId, accountId, input)
}

func (s *AccountService) GetGeneralAccount(userId, accountingId int) ([]models.CurrencyTotal, error) {
	_, err := s.accountingRepo.GetById(userId, accountingId)
	if err != nil {
		return nil, err
	}

	return s.repo.GetCurrencyTotals(accountingId)
}
