package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go-multi-tenancy/internals/core/domain"
	"go-multi-tenancy/internals/core/ports"
)

type companyService struct {
	companyRepository ports.CompanyRepository
}

func NewCompanyService(companyRepository ports.CompanyRepository) *companyService {
	return &companyService{
		companyRepository: companyRepository,
	}
}

func (s *companyService) Register(register *domain.RegisterInput) (*domain.DataReply, error) {
	if register.Username == "" || register.Password == "" || register.Company == "" || register.Branch == "" {
		return nil, errors.New("username or password cannot be empty")
	}

	hashedPassword := hashPassword(register.Password)
	registerData := &domain.Data{
		Company:   register.Company,
		Branch:    register.Branch,
		Username:  register.Username,
		Password:  hashedPassword,
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Role:      "user",
	}

	res, err := s.companyRepository.Register(registerData)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		ID:        res.ID,
		Company:   res.Company,
		Branch:    res.Branch,
		Username:  res.Username,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		CreatedAt: res.CreateAt,
	}, nil
}

func (s *companyService) Login(login *domain.LoginInput) (*domain.DataReply, string, error) {
	if login.Username == "" || login.Password == "" || login.Company == "" || login.Branch == "" {
		return nil, "", errors.New("username or password cannot be empty")
	}

	hashedPassword := hashPassword(login.Password)
	loginData := &domain.Data{
		Username: login.Username,
		Password: hashedPassword,
		Company:  login.Company,
		Branch:   login.Branch,
	}

	res, err := s.companyRepository.Login(loginData)
	if err != nil {
		return nil, "", err
	}

	if !checkPassword(login.Password, res.Password) {
		return nil, "", errors.New("invalid username or password")
	}

	return &domain.DataReply{
		ID:        res.ID,
		Company:   res.Company,
		Branch:    res.Branch,
		Username:  res.Username,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		CreatedAt: res.CreateAt,
	}, res.Role, nil
}

func (s *companyService) GetData(data *domain.DataInput) (*domain.DataReply, error) {
	req := &domain.Data{
		Company: data.Company,
		Branch:  data.Branch,
	}

	res, err := s.companyRepository.GetData(req)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		ID:        res.ID,
		Company:   res.Company,
		Branch:    res.Branch,
		Username:  res.Username,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		CreatedAt: res.CreateAt,
	}, nil
}

func (s *companyService) UpdateData(data *domain.DataUpdate) (*domain.DataReply, error) {
	req := &domain.Data{
		ID:        data.ID,
		Company:   data.Company,
		Branch:    data.Branch,
		Username:  data.Username,
		FirstName: *data.FirstName,
		LastName:  *data.LastName,
		Password:  *data.Password,
	}

	res, err := s.companyRepository.UpdateData(req)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		ID:        res.ID,
		Company:   res.Company,
		Branch:    res.Branch,
		Username:  res.Username,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		CreatedAt: res.CreateAt,
	}, nil
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func checkPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

func (s *companyService) DeleteData(data *domain.DataDelete) error {
	req := &domain.Data{
		Company: data.Company,
		Branch:  data.Branch,
		ID:      data.ID,
	}

	err := s.companyRepository.DeleteData(req)
	if err != nil {
		return err
	}

	return nil

}

func (s *companyService) GetAllData() ([]domain.DataReply, error) {
	res, err := s.companyRepository.GetAllData()
	if err != nil {
		return nil, err
	}
	data := []domain.DataReply{}
	for _, info := range res {
		data = append(data, domain.DataReply{
			ID:        info.ID,
			Company:   info.Company,
			Branch:    info.Branch,
			Username:  info.Username,
			FirstName: info.FirstName,
			LastName:  info.LastName,
			CreatedAt: info.CreateAt,
		})
	}

	return data, nil
}

func (s *companyService) GetMe(data *domain.Me) (*domain.DataReply, error) {
	form := &domain.Data{
		Company: data.Company,
		Branch:  data.Branch,
		ID:      data.ID,
	}

	res, err := s.companyRepository.GetMe(form)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		ID:        res.ID,
		Company:   res.Company,
		Branch:    res.Branch,
		Username:  res.Username,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		CreatedAt: res.CreateAt,
	}, nil
}

func (s *companyService) Admin(data *domain.Admin) (*domain.DataReply, error) {
	if data.Username == "" || data.Password == "" || data.Company == "" || data.Branch == "" {
		return nil, errors.New("username or password cannot be empty")
	}

	hashedPassword := hashPassword(data.Password)
	registerData := &domain.Data{
		Company:   data.Company,
		Branch:    data.Branch,
		Username:  data.Username,
		Password:  hashedPassword,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Role:      "admin",
	}

	res, err := s.companyRepository.Register(registerData)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		ID:        res.ID,
		Company:   res.Company,
		Branch:    res.Branch,
		Username:  res.Username,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		CreatedAt: res.CreateAt,
	}, nil

}

func (s *companyService) GetCompanyData(data *domain.DataInput) ([]domain.DataReply, error) {
	req := &domain.Data{
		Company: data.Company,
		Branch:  data.Branch,
	}

	res, err := s.companyRepository.GetCompanyData(req)
	if err != nil {
		return nil, err
	}

	dataCompany := []domain.DataReply{}
	for _, info := range res {
		dataCompany = append(dataCompany, domain.DataReply{
			ID:        info.ID,
			Company:   info.Company,
			Branch:    info.Branch,
			Username:  info.Username,
			FirstName: info.FirstName,
			LastName:  info.LastName,
			CreatedAt: info.CreateAt,
		})
	}

	return dataCompany, nil
}

func (s *companyService) GetBranchData(data *domain.DataInput) ([]domain.DataReply, error) {
	req := &domain.Data{
		Company: data.Company,
		Branch:  data.Branch,
	}

	res, err := s.companyRepository.GetBranchData(req)
	if err != nil {
		return nil, err
	}

	dataBranch := []domain.DataReply{}
	for _, info := range res {
		dataBranch = append(dataBranch, domain.DataReply{
			ID:        info.ID,
			Company:   info.Company,
			Branch:    info.Branch,
			Username:  info.Username,
			FirstName: info.FirstName,
			LastName:  info.LastName,
			CreatedAt: info.CreateAt,
		})
	}

	return dataBranch, nil
}
