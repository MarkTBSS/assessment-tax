package repository

import (
	"github.com/MarkTBSS/assessment-tax/databases"
	"github.com/MarkTBSS/assessment-tax/entities"
)

type deductionRepositoryImplement struct {
	db databases.Database
}

func NewTaxRepositoryImplement(db databases.Database) DeductionRepository {
	return &deductionRepositoryImplement{
		db: db,
	}
}

func (r *deductionRepositoryImplement) Listing() (*entities.Deduction, error) {
	query := r.db.Connect()
	deductions := entities.Deduction{}
	err := query.Find(&deductions).Error
	if err != nil {
		return nil, err
	}
	return &deductions, nil
}
