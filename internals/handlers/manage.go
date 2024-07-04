package handlers

import (
	"go-multi-tenancy/internals/core/domain"
	"go-multi-tenancy/internals/core/ports"

	"github.com/gofiber/fiber/v2"
)

type ManageHandler struct {
	manageService ports.ManageService
}

func NewManageHandler(manageService ports.ManageService) *ManageHandler {
	return &ManageHandler{manageService: manageService}
}

func (m *ManageHandler) GetCompany(c *fiber.Ctx) error {
	res, err := m.manageService.GetCompany()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": res,
	})
}

func (m *ManageHandler) GetBranch(c *fiber.Ctx) error {
	company := c.Params("company")

	req := &domain.CompanyRequest{
		Company: company,
	}

	res, err := m.manageService.GetBranch(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": res,
	})
}

func (m *ManageHandler) CreateCompany(c *fiber.Ctx) error {
	var company domain.CompanyRequest
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := m.manageService.CreateCompany(&company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": res.Company + " successfully created",
	})
}

func (m *ManageHandler) CreateBranch(c *fiber.Ctx) error {
	var branch domain.BranchRequest
	if err := c.BodyParser(&branch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := m.manageService.CreateBranch(&branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": *res.Branch + " in " + res.Company + " successfully created",
	})
}

func (m *ManageHandler) DeleteCompany(c *fiber.Ctx) error {
	company := c.Params("company")

	req := &domain.CompanyRequest{
		Company: company,
	}

	err := m.manageService.DeleteCompany(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": "Company successfully deleted",
	})
}

func (m *ManageHandler) DeleteBranch(c *fiber.Ctx) error {
	company := c.Params("company")
	branch := c.Params("branch")

	req := &domain.BranchRequest{
		Company: company,
		Branch:  branch,
	}

	err := m.manageService.DeleteBranch(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": "Branch successfully deleted",
	})
}

func (m *ManageHandler) UpdateCompanyToBranch(c *fiber.Ctx) error {
	company := c.Params("company")
	var req domain.CompanyAndBranch
	req.OldCompany = company
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := m.manageService.UpdateCompanyToBranch(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": *res.Branch + " Company successfully updated",
	})

}

func (m *ManageHandler) UpdateBranchToCompany(c *fiber.Ctx) error {
	branch := c.Params("branch")
	var req domain.CompanyAndBranch
	req.OldBranch = branch
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := m.manageService.UpdateBranchToCompany(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": res.Company + " Branch successfully updated",
	})

}

func (m *ManageHandler) UpdateCompanyName(c *fiber.Ctx) error {
	company := c.Params("company")
	var req domain.RenameCompany
	req.OldCompany = company
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := m.manageService.UpdateCompanyName(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": res.Company + " successfully updated",
	})
}

func (m *ManageHandler) UpdateBranchName(c *fiber.Ctx) error {
	branch := c.Params("branch")
	var req domain.RenameBranch
	req.OldBranch = branch
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := m.manageService.UpdateBranchName(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": *res.Branch + " successfully updated",
	})
}
