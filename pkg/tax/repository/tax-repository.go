package repository

import "github.com/MarkTBSS/assessment-tax/entities"

type DeductionRepository interface {
	Listing() (*entities.Deduction, error)
}
