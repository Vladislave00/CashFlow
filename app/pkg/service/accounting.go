package service

import (
	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/Vladislave00/CashFlow/app/pkg/repository"
)

type AccountingService struct {
	repo repository.Accounting
}

func NewAccountingService(repo repository.Accounting) *AccountingService {
	return &AccountingService{repo: repo}
}

func (s *AccountingService) CreateAccounting(userId int, accounting models.Accounting) (int, error) {
	return s.repo.CreateAccounting(userId, accounting)
}

func (s *AccountingService) GetAll(userId int) ([]models.Accounting, error) {
	return s.repo.GetAll(userId)
}

func (s *AccountingService) GetById(userId, accountingId int) (models.Accounting, error) {
	return s.repo.GetById(userId, accountingId)
}

func (s *AccountingService) Delete(userId, accountingId int) error {
	return s.repo.Delete(userId, accountingId)
}

func (s *AccountingService) Update(userId, accountingId int, input string) error {
	return s.repo.Update(userId, accountingId, input)
}
