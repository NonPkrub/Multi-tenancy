package ports

import (
	"go-multi-tenancy/internals/core/domain"

	"github.com/gofiber/fiber/v2"
)

type ManageRepository interface {
	GetCompany() ([]domain.Manage, error)
	GetBranch() ([]domain.Manage, error)
	CreateCompany(data *domain.Manage) (*domain.Manage, error)
	CreateBranch(data *domain.Manage) (*domain.Manage, error)
	DeleteCompany(data *domain.Manage) error
	DeleteBranch(data *domain.Manage) error
}

type ManageService interface {
	GetCompany() ([]domain.Response, error)
	GetBranch() ([]domain.Response, error)
	CreateCompany(data *domain.CompanyRequest) (*domain.Response, error)
	CreateBranch(data *domain.BranchRequest) (*domain.Response, error)
	DeleteCompany(data *domain.CompanyRequest) error
	DeleteBranch(data *domain.BranchRequest) error
}

type ManageHandler interface {
	GetCompany(c *fiber.Ctx) error
	GetBranch(c *fiber.Ctx) error
	CreateCompany(c *fiber.Ctx) error
	CreateBranch(c *fiber.Ctx) error
	DeleteCompany(c *fiber.Ctx) error
	DeleteBranch(c *fiber.Ctx) error
}
