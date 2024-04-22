package repository

import (
	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
)

type PersonalDeductionRepository interface {
	SettingPersonalDeduction(amount *model.AmountRequest) (*model.PersonalDeductionResponse, error)
	SettingKReceipt(amount *model.AmountRequest) (*model.KReceiptResponse, error)
}
