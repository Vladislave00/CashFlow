package service

import (
	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/Vladislave00/CashFlow/app/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Accounting interface {
	CreateAccounting(userId int, accounting models.Accounting) (int, error)
	GetAll(userId int) ([]models.Accounting, error)
	GetById(userId, accountingId int) (models.Accounting, error)
	Update(userId, accountingId int, input string) error
	Delete(userId, accountingId int) error
}

type Account interface {
	CreateAccount(userId, accountingId int, input models.CreateAccountInput) (int, error)
	GetAll(userId, accountingId int) ([]models.Account, error)
	GetById(userId, accountId int) (models.Account, error)
	Delete(userId, accountId int) error
	Update(userId, accountId int, input models.UpdateAccountInput) error
	GetGeneralAccount(userId, accountingId int) ([]models.CurrencyTotal, error)
}

type Transaction interface {
	CreateTransaction(userId, accountingId int, input models.Transaction) (int, error)
	GetAll(userId, accounting_id int) ([]models.Transaction, error)
	GetById(id, transaction_id int) (models.Transaction, error)
	Delete(userId, transaction_id int) error
	GetByAccountId(userId, accountingId, accountId int) ([]models.Transaction, error)
}

type Value interface {
	GetById(id int) (models.Money_value, error)
	GetByName(name string) (models.Money_value, error)
}

type Service struct {
	Authorization
	Accounting
	Account
	Transaction
	Value
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Accounting:    NewAccountingService(repo.Accounting),
		Value:         NewValueService(repo.Value),
		Account:       NewAccountService(repo.Account, repo.Accounting),
		Transaction:   NewTransactionService(repo.Transaction, repo.Account, repo.Accounting),
	}
}
