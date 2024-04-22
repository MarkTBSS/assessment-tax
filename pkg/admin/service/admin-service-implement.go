package service

import (
	"errors"

	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
	"github.com/MarkTBSS/assessment-tax/pkg/admin/repository"
)

type adminServiceImplement struct {
	adminRepository repository.PersonalDeductionRepository
}

func NewAdminServiceImplement(adminRepository repository.PersonalDeductionRepository) AdminService {
	return &adminServiceImplement{adminRepository}
}

func (s *adminServiceImplement) SetPersonalDeduction(amount *model.AmountRequest) (*model.PersonalDeductionResponse, error) {
	if amount.Amount < 10000.00 {
		return nil, errors.New("amount is required")
	}
	if amount.Amount > 100000.00 {
		amount.Amount = 100000.00
	}
	personalDeduction, err := s.adminRepository.Setting(amount)
	if err != nil {
		return nil, err
	}
	return &model.PersonalDeductionResponse{PersonalDeduction: personalDeduction.PersonalDeduction}, nil
}
