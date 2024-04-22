package controller

import "github.com/labstack/echo/v4"

type AdminController interface {
	SetPersonalDeduction(pctx echo.Context) error
	SetKReceipt(pctx echo.Context) error
}
