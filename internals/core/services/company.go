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
	if register.Username == "" || register.Password == "" || register.CompanyID == 0 || register.BranchID == 0 {
		return nil, errors.New("username or password cannot be empty")
	}

	hashedPassword := hashPassword(register.Password)
	registerData := &domain.Data{
		CompanyID: register.CompanyID,
		BranchID:  register.BranchID,
		Username:  register.Username,
		Password:  hashedPassword,
		DataValue: register.DataValue,
	}

	res, err := s.companyRepository.Register(registerData)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		CompanyID: res.CompanyID,
		BranchID:  res.BranchID,
		UserID:    res.UserID,
		Username:  res.Username,
		DataValue: res.DataValue,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (s *companyService) Login(login *domain.LoginInput) (*domain.DataReply, error) {
	if login.Username == "" || login.Password == "" || login.CompanyID == 0 || login.BranchID == 0 {
		return nil, errors.New("username or password cannot be empty")
	}

	hashedPassword := hashPassword(login.Password)
	loginData := &domain.Data{
		Username:  login.Username,
		Password:  hashedPassword,
		CompanyID: login.CompanyID,
		BranchID:  login.BranchID,
	}

	res, err := s.companyRepository.Login(loginData)
	if err != nil {
		return nil, err
	}

	if !checkPassword(login.Password, res.Password) {
		return nil, errors.New("invalid username or password")
	}

	return &domain.DataReply{
		CompanyID: res.CompanyID,
		BranchID:  res.BranchID,
		UserID:    res.UserID,
		Username:  res.Username,
		DataValue: res.DataValue,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (s *companyService) GetData(data *domain.DataInput) (*domain.DataReply, error) {
	req := &domain.Data{
		CompanyID: data.CompanyID,
		BranchID:  data.BranchID,
	}

	res, err := s.companyRepository.GetData(req)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		CompanyID: res.CompanyID,
		BranchID:  res.BranchID,
		UserID:    res.UserID,
		Username:  res.Username,
		DataValue: res.DataValue,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (s *companyService) UpdateData(data *domain.DataUpdate) (*domain.DataReply, error) {
	req := &domain.Data{
		CompanyID: data.CompanyID,
		BranchID:  data.BranchID,
		UserID:    data.UserID,
		DataValue: data.DataValue,
	}

	res, err := s.companyRepository.UpdateData(req)
	if err != nil {
		return nil, err
	}

	return &domain.DataReply{
		CompanyID: res.CompanyID,
		BranchID:  res.BranchID,
		UserID:    res.UserID,
		Username:  res.Username,
		DataValue: res.DataValue,
		CreatedAt: res.CreatedAt,
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
		CompanyID: data.CompanyID,
		BranchID:  data.BranchID,
		UserID:    data.UserID,
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
			CompanyID: info.CompanyID,
			BranchID:  info.BranchID,
			UserID:    info.UserID,
			Username:  info.Username,
			DataValue: info.DataValue,
			CreatedAt: info.CreatedAt,
		})
	}

	return data, nil
}
