package entities

type Deduction struct {
	PersonalDeduction float64 `gorm:"column:personal_deduction;not null"`
	KReceipt          float64 `gorm:"column:k_receipt;not null"`
}
