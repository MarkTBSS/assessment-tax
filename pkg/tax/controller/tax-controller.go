package controller

import "github.com/labstack/echo/v4"

type TaxController interface {
	Calculate(pctx echo.Context) error
	CalculateCSV(pctx echo.Context) error
}
