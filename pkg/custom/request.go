package custom

import (
	"encoding/csv"
	"io"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	EchoRequest interface {
		Bind(obj any) error
	}

	customEchoRequest struct {
		ctx       echo.Context
		validator *validator.Validate
	}
)

var (
	once              sync.Once
	validatorInstance *validator.Validate
)

func NewCustomEchoRequest(echoRequest echo.Context) EchoRequest {
	once.Do(func() {
		validatorInstance = validator.New()
	})
	return &customEchoRequest{
		ctx:       echoRequest,
		validator: validatorInstance,
	}
}

func (r *customEchoRequest) Bind(obj any) error {
	if err := r.ctx.Bind(obj); err != nil {
		return err
	}
	if err := r.validator.Struct(obj); err != nil {
		return err
	}
	return nil
}

// ParseCSV parses a CSV file and returns the records as a slice of maps.
func ParseCSV(file io.Reader) ([]map[string]float64, error) {
	reader := csv.NewReader(file)
	records := make([]map[string]float64, 0)
	// Read the header
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		record := make(map[string]float64)
		for i, value := range row {
			record[header[i]], err = strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
		}
		records = append(records, record)

	}
	return records, nil
}
