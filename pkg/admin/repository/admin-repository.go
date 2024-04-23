package repository

import (
	"github.com/MarkTBSS/assessment-tax/databases"
	"github.com/MarkTBSS/assessment-tax/entities"
	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
)

type DeductionRepository interface {
	SettingPersonalDeduction(amount *model.AmountRequest) (*model.PersonalDeductionResponse, error)
	SettingKReceipt(amount *model.AmountRequest) (*model.KReceiptResponse, error)
}

type DeductionRepositoryImplement struct {
	db databases.Database
}

func NewAdminRepositoryImplement(db databases.Database) DeductionRepository {
	return &DeductionRepositoryImplement{
		db: db,
	}
}

func (r *DeductionRepositoryImplement) SettingPersonalDeduction(request *model.AmountRequest) (*model.PersonalDeductionResponse, error) {
	personalDeductions := entities.Deduction{}
	err := r.db.Connect().Model(&personalDeductions).Table("deductions").Where("id = ?", 1).Update("personal_deduction", request.Amount).Error
	if err != nil {
		return nil, err
	}
	return &model.PersonalDeductionResponse{
		PersonalDeduction: personalDeductions.PersonalDeduction,
	}, nil
}

func (r *DeductionRepositoryImplement) SettingKReceipt(request *model.AmountRequest) (*model.KReceiptResponse, error) {
	kReceipt := entities.Deduction{}
	err := r.db.Connect().Model(&kReceipt).Table("deductions").Where("id = ?", 1).Update("k_receipt", request.Amount).Error
	if err != nil {
		return nil, err
	}
	return &model.KReceiptResponse{
		KReceipt: kReceipt.KReceipt,
	}, nil
}
