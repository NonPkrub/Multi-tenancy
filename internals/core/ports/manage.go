package ports

import (
	"go-multi-tenancy/internals/core/domain"

	"github.com/gofiber/fiber/v2"
)

type ManageRepository interface {
	GetCompany() ([]domain.Manage, error)
	GetBranch(data *domain.GetBranch) ([]domain.GetBranch, error)
	CreateCompany(data *domain.Manage) (*domain.Manage, error)
	CreateBranch(data *domain.Manage) (*domain.Manage, error)
	UpdateCompanyToBranch(data *domain.Manage) (*domain.Manage, error)
	UpdateBranchToCompany(data *domain.Manage) (*domain.Manage, error)
	UpdateCompanyName(data *domain.Manage) (*domain.Manage, error)
	UpdateBranchName(data *domain.Manage) (*domain.Manage, error)
	DeleteCompany(data *domain.Manage) error
	DeleteBranch(data *domain.Manage) error
}

type ManageService interface {
	GetCompany() ([]domain.ResponseCompany, error)
	GetBranch(data *domain.CompanyRequest) ([]domain.ResponseBranch, error)
	CreateCompany(data *domain.CompanyRequest) (*domain.Response, error)
	CreateBranch(data *domain.BranchRequest) (*domain.Response, error)
	UpdateCompanyToBranch(data *domain.BranchRequest) (*domain.Response, error)
	UpdateBranchToCompany(data *domain.BranchRequest) (*domain.Response, error)
	UpdateCompanyName(data *domain.BranchRequest) (*domain.Response, error)
	UpdateBranchName(data *domain.BranchRequest) (*domain.Response, error)
	DeleteCompany(data *domain.CompanyRequest) error
	DeleteBranch(data *domain.BranchRequest) error
}

type ManageHandler interface {
	GetCompany(c *fiber.Ctx) error
	GetBranch(c *fiber.Ctx) error
	CreateCompany(c *fiber.Ctx) error
	UpdateCompanyToBranch(c *fiber.Ctx) error
	UpdateBranchToCompany(c *fiber.Ctx) error
	UpdateCompanyName(c *fiber.Ctx) error
	UpdateBranchName(c *fiber.Ctx) error
	CreateBranch(c *fiber.Ctx) error
	DeleteCompany(c *fiber.Ctx) error
	DeleteBranch(c *fiber.Ctx) error
}
