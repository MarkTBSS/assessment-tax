package controller

import (
	"net/http"

	"github.com/MarkTBSS/assessment-tax/pkg/custom"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/model"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/service"
	"github.com/labstack/echo/v4"
)

type taxControllerImplement struct {
	taxService service.TaxService
}

func NewTaxControllerImplement(taxService service.TaxService) TaxController {
	return &taxControllerImplement{taxService}
}

func (c *taxControllerImplement) Calculate(pctx echo.Context) error {
	//var request model.TaxRequest
	request := new(model.TaxRequest)
	validatingContext := custom.NewCustomEchoRequest(pctx)
	if err := validatingContext.Bind(request); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	taxResult, err := c.taxService.Calculate(request)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}
	return pctx.JSON(http.StatusOK, taxResult)
}
