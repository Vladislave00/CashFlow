package repository

import (
	"fmt"

	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/jmoiron/sqlx"
)

type ValuePostgres struct {
	db *sqlx.DB
}

func NewValuePostgres(db *sqlx.DB) *ValuePostgres {
	return &ValuePostgres{db: db}
}

func (r *ValuePostgres) GetById(id int) (models.Money_value, error) {
	var value models.Money_value
	query := fmt.Sprintf("SELECT * FROM %s WHERE value_id=$1", valuesTable)
	err := r.db.Get(&value, query, id)

	return value, err
}

func (r *ValuePostgres) GetByName(name string) (models.Money_value, error) {
	var value models.Money_value
	query := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", valuesTable)
	err := r.db.Get(&value, query, name)

	return value, err
}
