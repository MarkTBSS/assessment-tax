package model

type AmountRequest struct {
	Amount float64 `json:"amount"`
}

type PersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}
