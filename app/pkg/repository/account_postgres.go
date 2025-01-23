package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/jmoiron/sqlx"
)

type AccountPostgres struct {
	db *sqlx.DB
}

func NewAccountPostgres(db *sqlx.DB) *AccountPostgres {
	return &AccountPostgres{db: db}
}

func (r *AccountPostgres) CreateAccount(accountingId int, account models.CreateAccountInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var value_id int

	queryForValueId := fmt.Sprintf("SELECT value_id FROM %s WHERE name=$1", valuesTable)
	if err := r.db.Get(&value_id, queryForValueId, account.Value); err != nil {
		return 0, err
	}

	money_amount, err := strconv.ParseFloat(account.MoneyAmount, 32)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (accounting_id, name, money_amount, value_id) values ($1, $2, $3, $4) RETURNING account_id", accountsTable)
	row := r.db.QueryRow(query, accountingId, account.Name, money_amount, value_id)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *AccountPostgres) GetAll(userId, accountingId int) ([]models.Account, error) {
	var accounts []models.Account
	query := fmt.Sprintf("SELECT t2.account_id, t2.accounting_id, t2.name, t2.money_amount, t2.value_id FROM %s t1 INNER JOIN %s t2 ON t1.accounting_id=t2.accounting_id WHERE t1.accounting_id=$1 AND t1.owner_id=$2", accountingsTable, accountsTable)
	if err := r.db.Select(&accounts, query, accountingId, userId); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountPostgres) GetById(userId, accountId int) (models.Account, error) {
	var account models.Account
	query := fmt.Sprintf("SELECT t2.account_id, t2.accounting_id, t2.name, t2.money_amount, t2.value_id FROM %s t1 INNER JOIN %s t2 ON t1.accounting_id=t2.accounting_id WHERE t2.account_id=$1 AND t1.owner_id=$2", accountingsTable, accountsTable)
	if err := r.db.Get(&account, query, accountId, userId); err != nil {
		return account, err
	}

	return account, nil
}

func (r *AccountPostgres) Delete(userId, accountId int) error {
	query := fmt.Sprintf("DELETE FROM %s t1 USING %s t2 WHERE t1.accounting_id=t2.accounting_id AND t2.owner_id=$1 AND t1.account_id=$2", accountsTable, accountingsTable)
	_, err := r.db.Exec(query, userId, accountId)
	return err
}

func (r *AccountPostgres) Update(userId, account_id int, input models.UpdateAccountInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Money_amount != nil {
		setValues = append(setValues, fmt.Sprintf("money_amount=$%d", argId))
		args = append(args, *input.Money_amount)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s t1 SET %s FROM %s t2
									WHERE t1.accounting_id = t2.accounting_id AND t2.owner_id = $%d AND t1.account_id = $%d`,
		accountsTable, setQuery, accountingsTable, argId, argId+1)
	args = append(args, userId, account_id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AccountPostgres) GetCurrencyTotals(accountingId int) ([]models.CurrencyTotal, error) {
	var totals []models.CurrencyTotal

	query := `
        SELECT 
            mv.name AS value, 
            SUM(a.money_amount) AS sum
        FROM accounts a
        JOIN money_values mv ON a.value_id = mv.value_id
        WHERE 
            a.accounting_id = $1
            AND a.name <> 'Внешний счёт'
        GROUP BY mv.name
        ORDER BY mv.name
    `

	err := r.db.Select(&totals, query, accountingId)
	if err != nil {
		return nil, fmt.Errorf("failed to get currency totals: %v", err)
	}

	return totals, nil
}
