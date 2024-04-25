package service

import (
	"errors"

	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
	"github.com/MarkTBSS/assessment-tax/pkg/admin/repository"
)

type AdminService interface {
	SetPersonalDeduction(amount *model.AmountRequest) (*model.PersonalDeductionResponse, error)
	SetKReceipt(amount *model.AmountRequest) (*model.KReceiptResponse, error)
}

type adminServiceImplement struct {
	adminRepository repository.DeductionRepository
}

func NewAdminServiceImplement(adminRepository repository.DeductionRepository) AdminService {
	return &adminServiceImplement{adminRepository}
}

func (s *adminServiceImplement) SetPersonalDeduction(amount *model.AmountRequest) (*model.PersonalDeductionResponse, error) {
	if amount.Amount < 10000.00 {
		return nil, errors.New("personal deduction must be greater than 10,000")
	}
	if amount.Amount > 100000.00 {
		return nil, errors.New("personal deduction must be lower than 100,000")
	}
	personalDeduction, err := s.adminRepository.SettingPersonalDeduction(amount)
	if err != nil {
		return nil, err
	}
	return &model.PersonalDeductionResponse{PersonalDeduction: personalDeduction.PersonalDeduction}, nil
}

func (s *adminServiceImplement) SetKReceipt(amount *model.AmountRequest) (*model.KReceiptResponse, error) {
	if amount.Amount < 0.00 {
		return nil, errors.New("k receipt must be greater than 0")
	}
	if amount.Amount > 100000.00 {
		return nil, errors.New("k receipt must be lower than 100,000")
	}
	kReceipt, err := s.adminRepository.SettingKReceipt(amount)
	if err != nil {
		return nil, err
	}
	return &model.KReceiptResponse{KReceipt: kReceipt.KReceipt}, nil
}
