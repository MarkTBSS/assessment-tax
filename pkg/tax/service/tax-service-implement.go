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
	deductionList, err := s.deductionRepository.Listing()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Personal Deduction : %f\n", deductionList.PersonalDeduction)
	fmt.Printf("K Reciept : %f\n", deductionList.KReceipt)
	fmt.Printf("Total Income : %f\n", taxRequest.TotalIncome)
	fmt.Printf("With Holding Tax : %f\n", taxRequest.WHT)
	fmt.Println(taxRequest.Allowances)
	fmt.Println(taxRequest.Allowances)

	totalIncome := taxRequest.TotalIncome

	// Deduction
	totalTaxableIncome := totalIncome - deductionList.PersonalDeduction

	for _, allowance := range taxRequest.Allowances {
		switch allowance.AllowanceType {
		case model.PersonalDeduction, model.KReceipt, model.Donation:
			totalTaxableIncome -= allowance.Amount
		}
	}

	fmt.Printf("Before Calculate : %f\n", totalTaxableIncome)

	/* การคำนวนภาษีตามขั้นบันใด
	รายได้ 0 - 150,000 ได้รับการยกเว้น
	150,001 - 500,000 อัตราภาษี 10%
	500,001 - 1,000,000 อัตราภาษี 15%
		1,000,001 - 2,000,000 อัตราภาษี 20%
	มากกว่า 2,000,000 อัตราภาษี 35% */

	tax := taxLevel(totalTaxableIncome)
	fmt.Printf("Tax : %f\n", tax)
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
