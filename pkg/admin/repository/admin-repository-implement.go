package repository

import (
	"github.com/MarkTBSS/assessment-tax/databases"
	"github.com/MarkTBSS/assessment-tax/entities"
	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
)

type personalDeductionRepositoryImplement struct {
	db databases.Database
}

func NewAdminRepositoryImplement(db databases.Database) PersonalDeductionRepository {
	return &personalDeductionRepositoryImplement{
		db: db,
	}
}

func (r *personalDeductionRepositoryImplement) Setting(request *model.AmountRequest) (*model.PersonalDeductionResponse, error) {
	personalDeductions := entities.PersonalDeduction{}
	err := r.db.Connect().Model(&personalDeductions).Table("deductions").Where("id = ?", 1).Update("personal_deduction", request.Amount).Error
	if err != nil {
		return nil, err
	}
	return &model.PersonalDeductionResponse{
		PersonalDeduction: personalDeductions.PersonalDeduction,
	}, nil
}
