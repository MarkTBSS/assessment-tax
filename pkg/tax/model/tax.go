package model

type AllowanceType string

const (
	// PersonalDeduction AllowanceType = "personal_deduction"
	Donation AllowanceType = "donation"
	KReceipt AllowanceType = "k-receipt"
)

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType AllowanceType `json:"allowanceType"`
	Amount        float64       `json:"amount"`
}

type TaxResponse struct {
	Tax       float64 `json:"tax,omitempty"`
	TaxRefund float64 `json:"taxRefund,omitempty"`
}

func IsAllowanceTypeCorrect(allowances []Allowance) bool {
	for _, allowance := range allowances {
		if allowance.AllowanceType != Donation && allowance.AllowanceType != KReceipt {
			return false
		}
	}
	return true
}
