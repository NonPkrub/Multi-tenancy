package handlers

import (
	"go-multi-tenancy/internals/core/domain"
	"go-multi-tenancy/internals/core/ports"
	"go-multi-tenancy/internals/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
)

var store = session.New()

type CompanyHandler struct {
	companyService ports.CompanyService
}

func NewCompanyHandler(companyService ports.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := h.companyService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, role, err := h.companyService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate JWT
	token, err := utils.GenerateJWT(res.Username, res.Company, res.Branch, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res, "token": token})
}

func (h *CompanyHandler) GetData(c *fiber.Ctx) error {
	company := c.Locals("company")
	branch := c.Locals("branch")

	if company == nil || branch == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	req := &domain.DataInput{
		Company: company.(string),
		Branch:  branch.(string),
	}

	res, err := h.companyService.GetData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) UpdateData(c *fiber.Ctx) error {
	company := c.Locals("company")
	branch := c.Locals("branch")
	id := c.Locals("id")

	if company == nil || branch == nil || id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	var dataValue domain.DataUpdate
	if err := c.BodyParser(&dataValue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	req := &domain.DataUpdate{
		Company: company.(string),
		Branch:  branch.(string),
		ID:      id.(uuid.UUID),
	}

	res, err := h.companyService.UpdateData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) GetAllData(c *fiber.Ctx) error {
	res, err := h.companyService.GetAllData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) DeleteData(c *fiber.Ctx) error {
	company := c.Locals("company_id")
	branch := c.Locals("branch_id")
	id := c.Locals("id")

	if company == nil || branch == nil || id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	req := &domain.DataDelete{
		Company: company.(string),
		Branch:  branch.(string),
		ID:      id.(uuid.UUID),
	}

	err := h.companyService.DeleteData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": "Data deleted successfully"})
}

func (h *CompanyHandler) GetMe(c *fiber.Ctx) error {
	company := c.Locals("company_id")
	branch := c.Locals("branch_id")
	id := c.Locals("id")

	if company == nil || branch == nil || id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	req := &domain.Me{
		Company: company.(string),
		Branch:  branch.(string),
		ID:      id.(uuid.UUID),
	}

	res, err := h.companyService.GetMe(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) Admin(c *fiber.Ctx) error {
	var req domain.Admin
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := h.companyService.Admin(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) GetCompanyData(c *fiber.Ctx) error {
	company := c.Params("company")

	req := &domain.DataInput{
		Company: company,
		Branch:  "",
	}

	res, err := h.companyService.GetCompanyData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) GetBranchData(c *fiber.Ctx) error {
	company := c.Params("company")
	branch := c.Params("branch")

	req := &domain.DataInput{
		Company: company,
		Branch:  branch,
	}

	res, err := h.companyService.GetBranchData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}
