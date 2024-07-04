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

func (m *manageService) GetCompany() ([]domain.ResponseCompany, error) {
	company, err := m.manageRepository.GetCompany()
	if err != nil {
		return nil, errors.New("Company not found")
	}

	companyData := []domain.ResponseCompany{}
	for _, res := range company {
		if len(res.Company) == 0 {
			return nil, errors.New("Company not found")
		}

		companyData = append(companyData, domain.ResponseCompany{
			Company: extractLastPart(res.Company),
		})
	}

	return companyData, nil
}

func (m *manageService) GetBranch(data *domain.CompanyRequest) ([]domain.ResponseBranch, error) {
	if data == nil || data.Company == "" {
		return nil, errors.New("Company name is required")
	}

	req := &domain.GetBranch{
		Company: strings.ToLower(data.Company),
	}

	branch, err := m.manageRepository.GetBranch(req)
	if err != nil {
		return nil, err
	}

	branchData := []domain.ResponseBranch{}
	for _, res := range branch {
		if len(res.Company) == 0 && len(res.Branch) == 0 {
			return nil, errors.New("Branch not found")
		}

		branchData = append(branchData, domain.ResponseBranch{
			Company: extractLastPart(res.Company),
			Branch:  extractLastPartFromSlice(res.Branch),
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

func extractLastPart(s string) string {
	parts := strings.Split(s, ".")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func extractLastPartFromSlice(slice []*string) []domain.BranchObject {
	var result []domain.BranchObject
	for _, str := range slice {
		if str != nil {
			extracted := extractLastPart(*str)
			result = append(result, domain.BranchObject{Name: extracted})
		}
	}
	return result
}

func (m *manageService) UpdateCompanyToBranch(data *domain.CompanyAndBranch) (*domain.Response, error) {
	if data.NewBranch == "" || data.NewCompany == "" || data.OldBranch == "" || data.OldCompany == "" || data.BranchName == "" {
		return nil, errors.New("All fields are required")
	}

	err := m.manageRepository.UpdateCompanyToBranch(data)
	if err != nil {
		return nil, err
	}

	return &domain.Response{
		Branch: &data.BranchName,
	}, nil
}

func (m *manageService) UpdateBranchToCompany(data *domain.CompanyAndBranch) (*domain.Response, error) {
	if data.NewBranch == "" || data.NewCompany == "" || data.OldBranch == "" || data.OldCompany == "" || data.BranchName == "" {
		return nil, errors.New("All fields are required")
	}

	err := m.manageRepository.UpdateBranchToCompany(data)
	if err != nil {
		return nil, err
	}

	return &domain.Response{
		Company: data.NewCompany,
	}, nil
}

func (m *manageService) UpdateCompanyName(data *domain.RenameCompany) (*domain.Response, error) {
	if data.OldCompany == "" || data.NewCompany == "" {
		return nil, errors.New("All fields are required")
	}

	err := m.manageRepository.UpdateCompanyName(data)
	if err != nil {
		return nil, err
	}

	return &domain.Response{
		Company: data.NewCompany,
	}, nil

}

func (m *manageService) UpdateBranchName(data *domain.RenameBranch) (*domain.Response, error) {
	if data.OldBranch == "" || data.NewBranch == "" || data.Company == "" {
		return nil, errors.New("All fields are required")
	}

	err := m.manageRepository.UpdateBranchName(data)
	if err != nil {
		return nil, err
	}

	return &domain.Response{
		Company: data.Company,
		Branch:  &data.NewBranch,
	}, nil

}
