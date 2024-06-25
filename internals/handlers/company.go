package handlers

import (
	"go-multi-tenancy/internals/core/domain"
	"go-multi-tenancy/internals/core/ports"
	"go-multi-tenancy/internals/utils"
	"strconv"

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

	res, err := h.companyService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate JWT
	token, err := utils.GenerateJWT(res.UserID, res.CompanyID, res.BranchID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res, "token": token})
}

func (h *CompanyHandler) GetData(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("companyID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
	}

	branchID, err := strconv.Atoi(c.Params("branchID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid branch ID"})
	}

	req := &domain.DataInput{
		CompanyID: companyID,
		BranchID:  branchID,
	}

	res, err := h.companyService.GetData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": res})
}

func (h *CompanyHandler) UpdateData(c *fiber.Ctx) error {
	companyID, err := strconv.Atoi(c.Params("companyID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
	}

	branchID, err := strconv.Atoi(c.Params("branchID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid branch ID"})
	}

	userID, err := strconv.Atoi(c.Params("userID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	var dataValue domain.DataUpdate
	if err := c.BodyParser(&dataValue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	req := &domain.DataUpdate{
		CompanyID: companyID,
		BranchID:  branchID,
		UserID:    userID,
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
	companyID, err := strconv.Atoi(c.Params("companyID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
	}

	branchID, err := strconv.Atoi(c.Params("branchID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid branch ID"})
	}

	userID, err := strconv.Atoi(c.Params("userID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid record ID"})
	}

	req := &domain.DataDelete{
		CompanyID: companyID,
		BranchID:  branchID,
		UserID:    userID,
	}

	err = h.companyService.DeleteData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"data": "Data deleted successfully"})
}
