package ports

import (
	"go-multi-tenancy/internals/core/domain"

	"github.com/gofiber/fiber/v2"
)

type CompanyService interface {
	Register(register *domain.RegisterInput) (*domain.DataReply, error)
	Login(login *domain.LoginInput) (*domain.DataReply, error)
	GetData(data *domain.DataInput) (*domain.DataReply, error)
	UpdateData(data *domain.DataUpdate) (*domain.DataReply, error)
	GetAllData() ([]domain.DataReply, error)
	DeleteData(data *domain.DataDelete) error
}

type CompanyRepository interface {
	Register(data *domain.Data) (*domain.Data, error)
	Login(data *domain.Data) (*domain.Data, error)
	GetData(data *domain.Data) (*domain.Data, error)
	UpdateData(data *domain.Data) (*domain.Data, error)
	GetAllData() ([]domain.Data, error)
	DeleteData(data *domain.Data) error
}

type CompanyHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetData(c *fiber.Ctx) error
	UpdateData(c *fiber.Ctx) error
	GetAllData(c *fiber.Ctx) error
	DeleteData(c *fiber.Ctx) error
}
