package server

import (
	"github.com/MarkTBSS/assessment-tax/pkg/admin/controller"
	"github.com/MarkTBSS/assessment-tax/pkg/admin/repository"
	"github.com/MarkTBSS/assessment-tax/pkg/admin/service"
)

func (s *echoServer) initAdminRouter() {
	router := s.app.Group("/admin/deductions")

	adminRepository := repository.NewAdminRepositoryImplement(s.db)
	adminService := service.NewAdminServiceImplement(adminRepository)
	adminController := controller.NewAdminControllerImplement(adminService)

	router.POST("/personal", adminController.SetPersonalDeduction)
}
