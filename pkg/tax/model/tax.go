package model

type AllowanceType string

const (
	Donation          AllowanceType = "donation"
	PersonalDeduction AllowanceType = "personal_deduction"
	KReceipt          AllowanceType = "k_receipt"
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
	Tax float64 `json:"tax"`
}
