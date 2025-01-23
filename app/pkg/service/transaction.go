package service

import (
	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/Vladislave00/CashFlow/app/pkg/repository"
)

type TransactionService struct {
	repo           repository.Transaction
	accountRepo    repository.Account
	accountingRepo repository.Accounting
}

func NewTransactionService(repo repository.Transaction, accountRepo repository.Account, accountingRepo repository.Accounting) *TransactionService {
	return &TransactionService{repo: repo, accountRepo: accountRepo, accountingRepo: accountingRepo}
}

func (s *TransactionService) CreateTransaction(userId, accountingId int, transaction models.Transaction) (int, error) {
	_, err := s.accountingRepo.GetById(userId, accountingId)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateTransaction(userId, accountingId, transaction)
}

func (s *TransactionService) GetAll(userId, accounting_id int) ([]models.Transaction, error) {
	_, err := s.accountingRepo.GetById(userId, accounting_id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(accounting_id)
}

func (s *TransactionService) GetByAccountId(userId, accountingId, accountId int) ([]models.Transaction, error) {
	_, err := s.accountingRepo.GetById(userId, accountingId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByAccountId(accountingId, accountId)
}

func (s *TransactionService) GetById(id, transaction_id int) (models.Transaction, error) {
	return s.repo.GetById(id, transaction_id)
}

func (s *TransactionService) Delete(userId, transaction_id int) error {
	return s.repo.Delete(userId, transaction_id)
}
