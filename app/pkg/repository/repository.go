package repository

import (
	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (models.User, error)
}

type Accounting interface {
	CreateAccounting(userId int, accounting models.Accounting) (int, error)
	GetAll(userId int) ([]models.Accounting, error)
	GetById(userId, accountingId int) (models.Accounting, error)
	Update(userId, accountingId int, input string) error
	Delete(userId, accountingId int) error
}

type Account interface {
	CreateAccount(accountingId int, input models.CreateAccountInput) (int, error)
	GetAll(userId, accountingId int) ([]models.Account, error)
	GetById(userId, accountId int) (models.Account, error)
	Delete(userId, accountId int) error
	Update(id, account_id int, input models.UpdateAccountInput) error
	GetCurrencyTotals(accountingId int) ([]models.CurrencyTotal, error)
}

type Transaction interface {
	CreateTransaction(userId, accountingId int, input models.Transaction) (int, error)
	GetAll(accounting_id int) ([]models.Transaction, error)
	GetById(id, transaction_id int) (models.Transaction, error)
	Delete(userId, transaction_id int) error
	GetByAccountId(accountingId, accountId int) ([]models.Transaction, error)
}

type Value interface {
	GetById(id int) (models.Money_value, error)
	GetByName(name string) (models.Money_value, error)
}

type Repository struct {
	Authorization
	Accounting
	Account
	Transaction
	Value
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Accounting:    NewAccountingPostgres(db),
		Value:         NewValuePostgres(db),
		Account:       NewAccountPostgres(db),
		Transaction:   NewTransactionPostgres(db),
	}
}
