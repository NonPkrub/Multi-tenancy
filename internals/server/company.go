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
	company.Post("/admin", s.company.Admin)
	company.Post("/login", s.company.Login)
	company.Use(middleware.JWTAuth())
	{
		company.Get("", middleware.AuthorizeRole("admin"), s.company.GetAllData)
		company.Get("/data", middleware.AuthorizeRole("admin"), s.company.GetData)
		company.Get("/data", s.company.GetMe)
		company.Put("/data", s.company.UpdateData)
		company.Delete("/data", s.company.DeleteData)
	}
	app.Listen(fmt.Sprintf(":%v", viper.GetInt("app.port")))

}
