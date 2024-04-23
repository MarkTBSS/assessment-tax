package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MarkTBSS/assessment-tax/pkg/custom"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/model"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/service"
	"github.com/labstack/echo/v4"
)

type TaxController interface {
	Calculate(pctx echo.Context) error
	CalculateCSV(pctx echo.Context) error
}

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

func (c *taxControllerImplement) CalculateCSV(pctx echo.Context) error {
	file, err := pctx.FormFile("taxFile")
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	if !strings.HasSuffix(file.Filename, ".csv") {
		return custom.Error(pctx, http.StatusBadRequest, errors.New("invalid file format, must be .csv"))
	}
	csvFile, err := file.Open()
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}
	defer csvFile.Close()
	records, err := custom.ParseCSV(csvFile)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	var taxResults []*model.TaxResponseCSV
	for _, record := range records {
		request := &model.TaxRequest{
			TotalIncome: record["totalIncome"],
			WHT:         record["wht"],
			Allowances: []model.Allowance{
				{
					AllowanceType: model.Donation,
					Amount:        record["donation"],
				},
			},
		}
		taxResult, err := c.taxService.Calculate(request)
		if err != nil {
			return custom.Error(pctx, http.StatusInternalServerError, err)
		}
		taxResults = append(taxResults, &model.TaxResponseCSV{
			TotalIncome: request.TotalIncome,
			Tax:         taxResult.Tax,
			TaxRefund:   taxResult.TaxRefund,
		})
	}
	return pctx.JSON(http.StatusOK, map[string]interface{}{"taxes": taxResults})
}
