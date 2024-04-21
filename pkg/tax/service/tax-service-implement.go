package service

import (
	"fmt"

	"github.com/MarkTBSS/assessment-tax/pkg/tax/model"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/repository"
)

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

	// Deduction Personal Deduction from Database
	totalIncome -= deductionList.PersonalDeduction
	fmt.Printf("After Personal Deduction : %f\n", totalIncome)

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
	fmt.Printf("After Allowance Deduction : %f\n", totalIncome)

	tax := taxLevel(totalIncome)
	fmt.Printf("Tax Before With Holding Tax Deduction : %f\n", tax)

	tax -= taxRequest.WHT
	fmt.Printf("Tax After With Holding Tax Deduction : %f\n", tax)

	if tax < 0 {
		return &model.TaxResponse{TaxRefund: -tax}, nil
	}
	return &model.TaxResponse{Tax: tax}, nil
}

func taxLevel(taxIncome float64) float64 {
	var tax float64
	tax = 0.0
	switch {
	case taxIncome <= 150000:
		tax = 0
	case taxIncome <= 500000:
		tax = (taxIncome - 150000) * 0.1
	case taxIncome <= 1000000:
		tax = 35000 + (taxIncome-500000)*0.15
	case taxIncome <= 2000000:
		tax = 100000 + (taxIncome-1000000)*0.2
	default:
		tax = 300000 + (taxIncome-2000000)*0.35
	}
	return tax
}
