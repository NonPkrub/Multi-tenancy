package handlers

import (
	"go-multi-tenancy/internals/core/domain"
	"go-multi-tenancy/internals/core/ports"
	"go-multi-tenancy/internals/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	token, err := utils.GenerateJWT(res.UserID, res.CompanyID, res.BranchID, role)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res, "token": token})
}

func (h *CompanyHandler) GetData(c *fiber.Ctx) error {
	companyID := c.Locals("company_id")
	branchID := c.Locals("branch_id")

	if companyID == nil || branchID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	req := &domain.DataInput{
		CompanyID: companyID.(int),
		BranchID:  branchID.(int),
	}

	res, err := h.companyService.GetData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) UpdateData(c *fiber.Ctx) error {
	companyID := c.Locals("company_id")
	branchID := c.Locals("branch_id")
	userID := c.Locals("user_id")

	if companyID == nil || branchID == nil || userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	var dataValue domain.DataUpdate
	if err := c.BodyParser(&dataValue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	req := &domain.DataUpdate{
		CompanyID: companyID.(int),
		BranchID:  branchID.(int),
		UserID:    userID.(int),
		DataValue: dataValue.DataValue,
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
	companyID := c.Locals("company_id")
	branchID := c.Locals("branch_id")
	userID := c.Locals("user_id")

	if companyID == nil || branchID == nil || userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	req := &domain.DataDelete{
		CompanyID: companyID.(int),
		BranchID:  branchID.(int),
		UserID:    userID.(int),
	}

	err := h.companyService.DeleteData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": "Data deleted successfully"})
}

func (h *CompanyHandler) GetMe(c *fiber.Ctx) error {
	companyID := c.Locals("company_id")
	branchID := c.Locals("branch_id")
	userID := c.Locals("user_id")

	if companyID == nil || branchID == nil || userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	req := &domain.Me{
		CompanyID: companyID.(int),
		BranchID:  branchID.(int),
		UserID:    userID.(int),
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
