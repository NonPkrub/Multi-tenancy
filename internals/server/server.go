package server

import (
	"fmt"
	"go-multi-tenancy/internals/core/ports"
	"go-multi-tenancy/internals/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

type Server struct {
	company ports.CompanyHandler
	manage  ports.ManageHandler
}

func NewServer(company ports.CompanyHandler, manage ports.ManageHandler) *Server {
	return &Server{company: company, manage: manage}
}

func (s *Server) Initialize() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	v1 := app.Group("/api/v1")
	company := v1.Group("company")
	company.Post("/register", s.company.Register) // create user
	company.Post("/admin", s.company.Admin)       // create admin
	company.Post("/login", s.company.Login)
	company.Use(middleware.JWTAuth())
	{
		company.Get("", middleware.AuthorizeRole("admin"), s.company.GetAllData)   // require admin role
		company.Get("/data", middleware.AuthorizeRole("admin"), s.company.GetData) // require admin role
		company.Get("/data", s.company.GetMe)
		company.Put("/data", s.company.UpdateData)
		company.Delete("/data", s.company.DeleteData)
	}

	manage := v1.Group("manage")
	{
		manage.Get("/company", s.manage.GetCompany)
		manage.Get("/branch/:company", s.manage.GetBranch)
		manage.Post("/company", s.manage.CreateCompany)
		manage.Post("/branch", s.manage.CreateBranch)
		manage.Put("/company/:company", s.manage.UpdateCompanyToBranch)
		manage.Put("/branch/:branch", s.manage.UpdateBranchToCompany)
		manage.Put("/rename/company/:company", s.manage.UpdateCompanyName)
		manage.Put("/rename/branch/:branch", s.manage.UpdateBranchName)
		manage.Delete("/company/:company", s.manage.DeleteCompany)
		manage.Delete("/company/:company/branch/:branch", s.manage.DeleteBranch)
	}

	app.Listen(fmt.Sprintf(":%v", viper.GetInt("app.port")))

}
