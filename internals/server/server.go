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
}

func NewServer(company ports.CompanyHandler) *Server {
	return &Server{company: company}
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
	app.Listen(fmt.Sprintf(":%v", viper.GetInt("app.port")))

}
