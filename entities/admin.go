package entities

type PersonalDeduction struct {
	PersonalDeduction float64 `gorm:"column:personal_deduction;not null"`
}

type KReceipt struct {
	KReceipt float64 `gorm:"column:k_receipt;not null"`
}
