package service

import (
	"github.com/Vladislave00/CashFlow/app/models"
	"github.com/Vladislave00/CashFlow/app/pkg/repository"
)

type ValueService struct {
	repo repository.Value
}

func NewValueService(repo repository.Value) *ValueService {
	return &ValueService{repo: repo}
}

func (s *ValueService) GetById(id int) (models.Money_value, error) {
	return s.repo.GetById(id)
}

func (s *ValueService) GetByName(name string) (models.Money_value, error) {
	return s.repo.GetByName(name)
}
