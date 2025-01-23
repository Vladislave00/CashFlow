package repository

import (
	"fmt"

	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/jmoiron/sqlx"
)

type AccountingPostgres struct {
	db *sqlx.DB
}

func NewAccountingPostgres(db *sqlx.DB) *AccountingPostgres {
	return &AccountingPostgres{db: db}
}

func (r *AccountingPostgres) CreateAccounting(userId int, accounting models.Accounting) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, owner_id) values ($1, $2) RETURNING accounting_id", accountingsTable)
	row := r.db.QueryRow(query, accounting.Name, userId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *AccountingPostgres) GetAll(userId int) ([]models.Accounting, error) {
	var accountings []models.Accounting
	query := fmt.Sprintf("SELECT * FROM %s WHERE owner_id=$1", accountingsTable)
	err := r.db.Select(&accountings, query, userId)

	return accountings, err
}

func (r *AccountingPostgres) GetById(userId, accountingId int) (models.Accounting, error) {
	var accounting models.Accounting
	query := fmt.Sprintf("SELECT * FROM %s WHERE owner_id=$1 AND accounting_id=$2", accountingsTable)
	err := r.db.Get(&accounting, query, userId, accountingId)

	return accounting, err
}

func (r *AccountingPostgres) Delete(userId, accountingId int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE owner_id=$1 AND accounting_id=$2", accountingsTable)
	_, err := r.db.Exec(query, userId, accountingId)

	return err
}

func (r *AccountingPostgres) Update(userId, accountingId int, input string) error {

	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE owner_id=$2 AND accounting_id=$3", accountingsTable)

	_, err := r.db.Exec(query, input, userId, accountingId)
	return err
}
