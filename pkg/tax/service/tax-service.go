package service

import (
	"fmt"

	"github.com/MarkTBSS/assessment-tax/entities"
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
	if taxRequest.TotalIncome < taxRequest.WHT || taxRequest.WHT < 0 {
		return nil, fmt.Errorf("with holding tax must be lower than total income and greater than 0")
	}
	deductionList, err := s.deductionRepository.Listing()
	if err != nil {
		return nil, err
	}
	totalIncome := calculateTotalIncome(taxRequest, deductionList)
	tax, taxLevels := TaxLevel(totalIncome)
	tax -= taxRequest.WHT
	if tax < 0 {
		return &model.TaxResponse{TaxRefund: -tax, TaxLevel: taxLevels}, nil
	} else {
		return &model.TaxResponse{Tax: tax, TaxLevel: taxLevels}, nil
	}
}

func calculateTotalIncome(taxRequest *model.TaxRequest, deductionList *entities.Deduction) float64 {
	totalIncome := taxRequest.TotalIncome - deductionList.PersonalDeduction
	for i := range taxRequest.Allowances {
		switch taxRequest.Allowances[i].AllowanceType {
		case model.Donation:
			if taxRequest.Allowances[i].Amount > 100000 {
				taxRequest.Allowances[i].Amount = 100000
			}
			totalIncome -= taxRequest.Allowances[i].Amount
		case model.KReceipt:
			if taxRequest.Allowances[i].Amount > deductionList.KReceipt || taxRequest.Allowances[i].Amount < 0 {
				taxRequest.Allowances[i].Amount = deductionList.KReceipt
			}
			totalIncome -= taxRequest.Allowances[i].Amount
		default:
			return 0
		}
	}
	return totalIncome
}

func TaxLevel(taxIncome float64) (float64, []model.TaxLevel) {
	tax, totalTax := 0.0, 0.0
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

/* func TaxLevel(taxIncome float64) (float64, []model.TaxLevel) {
	var taxLevels []model.TaxLevel

	taxBands := []struct {
		minIncome float64
		maxIncome float64
		rate      float64
		baseTax   float64
		level     string
	}{
		{0.0, 150000.0, 0.0, 0.0, "0-150,000"},
		{150000.0, 500000.0, 0.1, 0.0, "150,001-500,000"},
		{500000.0, 1000000.0, 0.15, 35000.0, "500,001-1,000,000"},
		{1000000.0, 2000000.0, 0.2, 135000.0, "1,000,001-2,000,000"},
		{2000000.0, math.MaxFloat64, 0.35, 435000.0, "2,000,001 ขึ้นไป"},
	}

	totalTax := 0.0
	for _, band := range taxBands {
		if taxIncome > band.minIncome {
			taxableIncome := math.Min(band.maxIncome, taxIncome) - band.minIncome
			fmt.Println("=====")
			fmt.Println("band.maxIncome", band.maxIncome)
			fmt.Println("band.minIncome", band.minIncome)
			fmt.Println("taxableIncome", taxableIncome)
			tax := taxableIncome * band.rate
			fmt.Println("tax", tax)
			totalTax += tax
			fmt.Println("totalTax", totalTax)
			taxLevels = append(taxLevels, model.TaxLevel{Level: band.level, Tax: band.baseTax + tax})
		}
	}
	return totalTax, taxLevels
} */
