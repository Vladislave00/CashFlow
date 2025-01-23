package repository

import (
	"fmt"

	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/jmoiron/sqlx"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) CreateTransaction(userId, accountingId int, input models.Transaction) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (account_id, external_account_id, money_amount) values ($1, $2, $3) RETURNING transaction_id", transactionsTable)
	row := r.db.QueryRow(query, input.Account_id, input.External_account_id, input.Money_amount)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TransactionPostgres) GetAll(accounting_id int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := fmt.Sprintf("SELECT t2.transaction_id, t2.account_id, t2.external_account_id, t2.money_amount, t2.created_at FROM %s t1 INNER JOIN %s t2 ON t1.account_id=t2.account_id WHERE t1.accounting_id=$1", accountsTable, transactionsTable)
	if err := r.db.Select(&transactions, query, accounting_id); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionPostgres) GetById(user_id, transaction_id int) (models.Transaction, error) {
	var transaction models.Transaction
	query := fmt.Sprintf("SELECT t1.transaction_id, t1.account_id, t1.external_account_id, t1.money_amount, t1.created_at FROM %s t1 INNER JOIN %s t2 USING(account_id) INNER JOIN %s t3 USING(accounting_id) WHERE t1.transaction_id=$1 AND t3.owner_id=$2", transactionsTable, accountsTable, accountingsTable)
	if err := r.db.Get(&transaction, query, transaction_id, user_id); err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *TransactionPostgres) Delete(userId, transaction_id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE transaction_id = $1 AND account_id IN (SELECT account_id FROM %s WHERE accounting_id IN (SELECT accounting_id FROM %s WHERE owner_id = $2))", transactionsTable, accountsTable, accountingsTable)
	_, err := r.db.Exec(query, transaction_id, userId)

	return err
}

func (r *TransactionPostgres) GetByAccountId(accountingId, accountId int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := fmt.Sprintf("SELECT t2.transaction_id, t2.account_id, t2.external_account_id, t2.money_amount, t2.created_at FROM %s t1 INNER JOIN %s t2 ON t1.account_id=t2.account_id WHERE t1.accounting_id=$1 AND t2.account_id=$2 OR t2.external_account_id=$2", accountsTable, transactionsTable)
	if err := r.db.Select(&transactions, query, accountingId, accountId); err != nil {
		return nil, err
	}

	return transactions, nil
}
