package service

import (
	"github.com/MarkTBSS/assessment-tax/pkg/tax/model"
)

type TaxService interface {
	Calculate(taxRequest *model.TaxRequest) (*model.TaxResponse, error)
}
