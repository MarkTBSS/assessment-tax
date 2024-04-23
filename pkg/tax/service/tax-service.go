package service

import (
	"fmt"

	"github.com/MarkTBSS/assessment-tax/pkg/tax/model"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/repository"
)

type TaxService interface {
	Calculate(taxRequest *model.TaxRequest) (*model.TaxResponse, error)
}

type taxServiceImplement struct {
	deductionRepository repository.DeductionRepository
}

func NewTaxServiceImplement(deductionRepository repository.DeductionRepository) TaxService {
	return &taxServiceImplement{deductionRepository}
}
func (s *taxServiceImplement) Calculate(taxRequest *model.TaxRequest) (*model.TaxResponse, error) {
	if !model.IsAllowanceTypeCorrect(taxRequest.Allowances) {
		return nil, fmt.Errorf("incorrect allowance type")
	}
	deductionList, err := s.deductionRepository.Listing()
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Personal Deduction : %f\n", deductionList.PersonalDeduction)
	//fmt.Printf("K Reciept : %f\n", deductionList.KReceipt)
	//fmt.Printf("Total Income : %f\n", taxRequest.TotalIncome)
	//fmt.Printf("With Holding Tax : %f\n", taxRequest.WHT)
	//fmt.Println(taxRequest.Allowances)

	totalIncome := taxRequest.TotalIncome
	totalIncome -= deductionList.PersonalDeduction
	//fmt.Printf("After Personal Deduction : %f\n", totalIncome)
	for _, allowance := range taxRequest.Allowances {
		switch allowance.AllowanceType {
		case model.Donation:
			if allowance.Amount > 100000 {
				allowance.Amount = 100000
			}
			totalIncome -= allowance.Amount
		case model.KReceipt:
			if allowance.Amount > deductionList.KReceipt {
				allowance.Amount = deductionList.KReceipt
			}
			totalIncome -= allowance.Amount
		}
	}
	//fmt.Printf("After Allowance Deduction : %f\n", totalIncome)
	tax, taxLevels := taxLevel(totalIncome)
	//fmt.Printf("Tax Before With Holding Tax Deduction : %f\n", tax)
	tax -= taxRequest.WHT
	//fmt.Printf("Tax After With Holding Tax Deduction : %f\n", tax)
	if tax < 0 {
		return &model.TaxResponse{TaxRefund: -tax, TaxLevel: taxLevels}, nil
	} else {
		return &model.TaxResponse{Tax: tax, TaxLevel: taxLevels}, nil
	}
}

func taxLevel(taxIncome float64) (float64, []model.TaxLevel) {
	var tax float64
	tax = 0.0
	var totalTax float64
	totalTax = 0.0
	var taxLevels []model.TaxLevel

	switch {
	case taxIncome <= 150000:
		taxLevels = append(taxLevels, model.TaxLevel{Level: "0-150,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "150,001-500,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "500,001-1,000,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "1,000,001-2,000,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: 0})
	case taxIncome <= 500000:
		tax = (taxIncome - 150000) * 0.1
		totalTax = tax
		taxLevels = append(taxLevels, model.TaxLevel{Level: "0-150,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "150,001-500,000", Tax: tax})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "500,001-1,000,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "1,000,001-2,000,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: 0})
	case taxIncome <= 1000000:
		tax = (taxIncome - 500000) * 0.15
		totalTax = 35000 + tax
		taxLevels = append(taxLevels, model.TaxLevel{Level: "0-150,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "150,001-500,000", Tax: 35000})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "500,001-1,000,000", Tax: tax})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "1,000,001-2,000,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: 0})
	case taxIncome <= 2000000:
		tax = (taxIncome - 1000000) * 0.2
		totalTax = 135000 + tax
		taxLevels = append(taxLevels, model.TaxLevel{Level: "0-150,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "150,001-500,000", Tax: 35000})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "500,001-1,000,000", Tax: 100000})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "1,000,001-2,000,000", Tax: tax})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: 0})
	default:
		tax = (taxIncome - 2000000) * 0.35
		totalTax = 435000 + tax
		taxLevels = append(taxLevels, model.TaxLevel{Level: "0-150,000", Tax: 0})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "150,001-500,000", Tax: 35000})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "500,001-1,000,000", Tax: 100000})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "1,000,001-2,000,000", Tax: 300000})
		taxLevels = append(taxLevels, model.TaxLevel{Level: "2,000,001 ขึ้นไป", Tax: tax})
	}
	return totalTax, taxLevels
}
