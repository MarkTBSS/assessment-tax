package service

import (
	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
)

type AdminService interface {
	SetPersonalDeduction(amount *model.AmountRequest) (*model.PersonalDeductionResponse, error)
}
