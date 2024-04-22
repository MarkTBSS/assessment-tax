package repository

import (
	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
)

type PersonalDeductionRepository interface {
	Setting(amount *model.AmountRequest) (*model.PersonalDeductionResponse, error)
}
