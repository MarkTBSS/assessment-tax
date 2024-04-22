package server

import (
	"github.com/MarkTBSS/assessment-tax/pkg/tax/controller"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/repository"
	"github.com/MarkTBSS/assessment-tax/pkg/tax/service"
)

func (s *echoServer) initTaxRouter() {
	router := s.app.Group("/tax/calculations")

	taxRepository := repository.NewTaxRepositoryImplement(s.db)
	taxService := service.NewTaxServiceImplement(taxRepository)
	taxController := controller.NewTaxControllerImplement(taxService)

	router.GET("", taxController.Calculate)
	router.POST("/upload-csv", taxController.CalculateCSV)
}
