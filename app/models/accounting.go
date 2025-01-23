package models

import "errors"

type Accounting struct {
	ID       int    `json:"id" db:"accounting_id"`
	Name     string `json:"name" db:"name" binding:"required"`
	Owner_id int    `json:"owner_id" db:"owner_id"`
}

type Money_value struct {
	ID   int     `json:"value_id" db:"value_id"`
	Name string  `json:"name" db:"name"`
	Rub  float32 `json:"rub" db:"rub"`
}

type Account struct {
	ID            int     `json:"account_id" db:"account_id"`
	Name          string  `json:"name" db:"name"`
	Accounting_id int     `json:"accounting_id" db:"accounting_id"`
	Money_amount  float32 `json:"money_amount" db:"money_amount"`
	Value_id      int     `json:"value_id" db:"value_id"`
}

type Transaction struct {
	ID                  int     `json:"transaction_id" db:"transaction_id"`
	Account_id          int     `json:"account_id" db:"account_id"`
	External_account_id int     `json:"external_account_id" db:"external_account_id"`
	Money_amount        float32 `json:"money_amount" db:"money_amount"`
	Created_at          string  `json:"created_at" db:"created_at"`
}

type UpdateAccountingInput struct {
	Name *string `json:"name"`
}

type CreateAccountInput struct {
	Name        string `json:"name"`
	MoneyAmount string `json:"money_amount"`
	Value       string `json:"value"`
}

type UpdateAccountInput struct {
	Name         *string  `json:"name"`
	Money_amount *float32 `json:"money_amount"`
}

func (i UpdateAccountInput) Validate() error {
	if i.Name == nil && i.Money_amount == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type CurrencyTotal struct {
	Value string  `json:"value"`
	Sum   float64 `json:"sum"`
}
