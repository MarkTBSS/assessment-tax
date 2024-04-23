package controller

import (
	"errors"
	"net/http"
	"os"

	"github.com/MarkTBSS/assessment-tax/pkg/admin/model"
	"github.com/MarkTBSS/assessment-tax/pkg/admin/service"
	"github.com/MarkTBSS/assessment-tax/pkg/custom"
	"github.com/labstack/echo/v4"
)

type AdminController interface {
	SetPersonalDeduction(pctx echo.Context) error
	SetKReceipt(pctx echo.Context) error
}

type adminControllerImplement struct {
	adminService service.AdminService
}

func NewAdminControllerImplement(adminService service.AdminService) AdminController {
	return &adminControllerImplement{adminService}
}

func (c *adminControllerImplement) SetPersonalDeduction(pctx echo.Context) error {
	if !checkBasicAuth(pctx) {
		return custom.Error(pctx, http.StatusUnauthorized, errors.New("unauthorized to set personal deduction"))
	}
	request := new(model.AmountRequest)
	validatingContext := custom.NewCustomEchoRequest(pctx)
	if err := validatingContext.Bind(request); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	adminResult, err := c.adminService.SetPersonalDeduction(request)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}
	return pctx.JSON(http.StatusOK, adminResult)
}

func (c *adminControllerImplement) SetKReceipt(pctx echo.Context) error {
	if !checkBasicAuth(pctx) {
		return custom.Error(pctx, http.StatusUnauthorized, errors.New("unauthorized to set k receipt"))
	}
	request := new(model.AmountRequest)
	validatingContext := custom.NewCustomEchoRequest(pctx)
	if err := validatingContext.Bind(request); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	adminResult, err := c.adminService.SetKReceipt(request)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}
	return pctx.JSON(http.StatusOK, adminResult)
}

func checkBasicAuth(ctx echo.Context) bool {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	user, pass, ok := ctx.Request().BasicAuth()
	if !ok {
		return false
	}
	return username == user && password == pass
}
