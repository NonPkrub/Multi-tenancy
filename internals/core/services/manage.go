package services

import (
	"errors"
	"go-multi-tenancy/internals/core/domain"
	"go-multi-tenancy/internals/core/ports"
	"strings"
)

type manageService struct {
	manageRepository ports.ManageRepository
}

func NewManageService(manageRepository ports.ManageRepository) *manageService {
	return &manageService{
		manageRepository: manageRepository,
	}
}

func (m *manageService) GetCompany() ([]domain.Response, error) {
	company, err := m.manageRepository.GetCompany()
	if err != nil {
		return nil, errors.New("Company not found")
	}

	companyData := []domain.Response{}
	for _, res := range company {
		if len(res.Company) == 0 {
			return nil, errors.New("Company not found")
		}

		companyData = append(companyData, domain.Response{
			Company: res.Company,
			Branch:  &res.Branch,
		})
	}

	return companyData, nil
}

func (m *manageService) GetBranch(data *domain.CompanyRequest) ([]domain.Response, error) {
	req := &domain.Manage{
		Company: strings.ToLower(data.Company),
	}

	branch, err := m.manageRepository.GetBranch(req)
	if err != nil {
		return nil, err
	}

	branchData := []domain.Response{}
	for _, res := range branch {
		if len(res.Company) == 0 && len(res.Branch) == 0 {
			return nil, errors.New("Branch not found")
		}

		branchData = append(branchData, domain.Response{
			Company: res.Company,
			Branch:  &res.Branch,
		})
	}

	return branchData, nil

}

func (m *manageService) CreateCompany(data *domain.CompanyRequest) (*domain.Response, error) {
	req := &domain.Manage{
		Company: strings.ToLower(data.Company),
	}

	company, err := m.manageRepository.CreateCompany(req)
	if err != nil {
		return nil, err
	}

	companyData := domain.Response{
		Company: company.Company,
		Branch:  &company.Branch,
	}

	return &companyData, nil
}

func (m *manageService) CreateBranch(data *domain.BranchRequest) (*domain.Response, error) {
	req := &domain.Manage{
		Company: strings.ToLower(data.Company),
		Branch:  strings.ToLower(data.Branch),
	}

	branch, err := m.manageRepository.CreateBranch(req)
	if err != nil {
		return nil, err
	}

	branchData := domain.Response{
		Company: branch.Company,
		Branch:  &branch.Branch,
	}

	return &branchData, nil
}

func (m *manageService) DeleteCompany(data *domain.CompanyRequest) error {
	req := &domain.Manage{
		Company: strings.ToLower(data.Company),
	}

	err := m.manageRepository.DeleteCompany(req)
	if err != nil {
		return err
	}

	return nil
}

func (m *manageService) DeleteBranch(data *domain.BranchRequest) error {
	req := &domain.Manage{
		Company: strings.ToLower(data.Company),
		Branch:  strings.ToLower(data.Branch),
	}

	err := m.manageRepository.DeleteBranch(req)
	if err != nil {
		return err
	}

	return nil
}
