package service_test

import (
	"testing"

	"github.com/MarkTBSS/assessment-tax/entities"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/model"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/service"
)

type MockDeductionRepository struct {
	PersonalDeduction float64
	KReceipt          float64
}

func (m *MockDeductionRepository) Listing() (*entities.Deduction, error) {
	return &entities.Deduction{
		PersonalDeduction: m.PersonalDeduction,
		KReceipt:          m.KReceipt,
	}, nil
}

func TestCalculateTax(t *testing.T) {
	mockDeductionList := &MockDeductionRepository{
		PersonalDeduction: 60000.0,
		KReceipt:          50000.0,
	}
	taxService := service.NewTaxServiceImplement(mockDeductionList)
	testCases := []struct {
		Name          string
		TaxRequest    *model.TaxRequest
		ExpectedTax   float64
		ExpectedLevel []model.TaxLevel
		ExpectedErr   error
	}{
		{
			Name: "Calculate Tax",
			TaxRequest: &model.TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []model.Allowance{
					{
						AllowanceType: model.KReceipt,
						Amount:        200000.0,
					},
					{
						AllowanceType: model.Donation,
						Amount:        100000.0,
					},
				},
			},
			ExpectedTax: 14000.0,
			ExpectedLevel: []model.TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   14000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   0.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
			ExpectedErr: nil,
		},
		{
			Name: "Calculate Tax + Donation more than 100,000",
			TaxRequest: &model.TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []model.Allowance{
					{
						AllowanceType: model.KReceipt,
						Amount:        200000.0,
					},
					{
						AllowanceType: model.Donation,
						Amount:        500000.0,
					},
				},
			},
			ExpectedTax: 14000.0,
			ExpectedLevel: []model.TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   14000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   0.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
			ExpectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			taxResult, err := taxService.Calculate(tc.TaxRequest)
			if err != tc.ExpectedErr {
				t.Errorf("Unexpected error: got %v, want %v", err, tc.ExpectedErr)
			}
			if taxResult.Tax != tc.ExpectedTax {
				t.Errorf("Incorrect Tax. Expected: %f, Got: %f", tc.ExpectedTax, taxResult.Tax)
			}
			if !areTaxLevelsEqual(taxResult.TaxLevel, tc.ExpectedLevel) {
				t.Errorf("Incorrect Tax Levels")
			}
		})
	}
}

func TestTaxLevel(t *testing.T) {
	testCases := []struct {
		Name           string
		TaxIncome      float64
		ExpectedTax    float64
		ExpectedLevels []model.TaxLevel
	}{
		{
			Name:        "0-150,000",
			TaxIncome:   150000.0,
			ExpectedTax: 0.0,
			ExpectedLevels: []model.TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   0.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   0.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		},
		{
			Name:        "150,001-500,000",
			TaxIncome:   300000.0,
			ExpectedTax: 15000.0,
			ExpectedLevels: []model.TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   15000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   0.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		},
		{
			Name:        "500,001-1,000,000",
			TaxIncome:   800000.0,
			ExpectedTax: 80000.0,
			ExpectedLevels: []model.TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   35000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   45000.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		},
		{
			Name:        "1,000,001-2,000,000",
			TaxIncome:   1500000.0,
			ExpectedTax: 235000.0,
			ExpectedLevels: []model.TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   35000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   100000.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   100000.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		},
		{
			Name:        "2,000,001 ขึ้นไป",
			TaxIncome:   2500000.0,
			ExpectedTax: 610000.0,
			ExpectedLevels: []model.TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   35000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   100000.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   300000.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   175000.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tax, taxLevels := service.TaxLevel(tc.TaxIncome)
			if tax != tc.ExpectedTax {
				t.Errorf("Incorrect Tax. Expected: %f, Got: %f", tc.ExpectedTax, tax)
			}
			if !areTaxLevelsEqual(taxLevels, tc.ExpectedLevels) {
				t.Errorf("Incorrect Tax Levels")
			}
		})
	}
}

func areTaxLevelsEqual(a, b []model.TaxLevel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Level != b[i].Level || a[i].Tax != b[i].Tax {
			return false
		}
	}
	return true
}
