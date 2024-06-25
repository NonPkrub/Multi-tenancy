package server

import (
	"fmt"
	"go-multi-tenancy/internals/core/ports"
	"go-multi-tenancy/internals/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type CompanyServer struct {
	company ports.CompanyHandler
}

func NewCompanyServer(company ports.CompanyHandler) *CompanyServer {
	return &CompanyServer{company: company}
}

func (s *CompanyServer) Initialize() {
	app := fiber.New()

	v1 := app.Group("/api/v1")

	company := v1.Group("company")
	company.Post("/register", s.company.Register)
	company.Post("/login", s.company.Login)
	company.Get("", s.company.GetAllData)
	company.Use(middleware.JWTAuth())
	{
		company.Get("/data/:companyID/:branchID", s.company.GetData)
		company.Put("/data/:companyID/:branchID/:userID", s.company.UpdateData)
		company.Delete("/data/:companyID/:branchID/:userID", s.company.DeleteData)
	}
	app.Listen(fmt.Sprintf(":%v", viper.GetInt("app.port")))

}
