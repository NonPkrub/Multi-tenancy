package repositories

import (
	"go-multi-tenancy/internals/core/domain"

	"github.com/jmoiron/sqlx"
)

type companyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *companyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Register(data *domain.Data) (*domain.Data, error) {
	query := "INSERT INTO company.test (username, password, company_id, branch_id, data_value, role)   VALUES ($1, $2, $3, $4, $5, $6) RETURNING company_id, branch_id, user_id, username, data_value, created_at"
	err := r.db.QueryRow(query, data.Username, data.Password, data.CompanyID, data.BranchID, data.DataValue, data.Role).Scan(&data.CompanyID, &data.BranchID, &data.UserID, &data.Username, &data.DataValue, &data.CreatedAt)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) Login(data *domain.Data) (*domain.Data, error) {
	query := "SELECT * FROM company.test WHERE username = $1 AND company_id = $2 AND branch_id = $3"
	err := r.db.QueryRow(query, data.Username, data.CompanyID, data.BranchID).Scan(&data.CompanyID, &data.BranchID, &data.UserID, &data.DataValue, &data.CreatedAt, &data.Username, &data.Password, &data.Role)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) GetData(data *domain.Data) (*domain.Data, error) {
	query := " SELECT * FROM company.test WHERE company_id = $1 AND branch_id = $2 "
	err := r.db.QueryRow(query, data.CompanyID, data.BranchID).Scan(&data.CompanyID, &data.BranchID, &data.UserID, &data.DataValue, &data.CreatedAt, &data.Username, &data.Password, &data.Role)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) UpdateData(data *domain.Data) (*domain.Data, error) {
	query := "UPDATE company.test SET data_value = $1 WHERE company_id = $2 AND branch_id = $3 AND user_id = $4 returning data_value, company_id, branch_id, user_id, username, created_at"
	err := r.db.QueryRow(query, data.DataValue, data.CompanyID, data.BranchID, data.UserID).Scan(&data.DataValue, &data.CompanyID, &data.BranchID, &data.UserID, &data.Username, &data.CreatedAt)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) DeleteData(data *domain.Data) error {
	query := "DELETE FROM company.test WHERE company_id = $1 AND branch_id = $2 AND user_id = $3"
	_, err := r.db.Exec(query, data.CompanyID, data.BranchID, data.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *companyRepository) GetAllData() ([]domain.Data, error) {
	data := []domain.Data{}

	query := "SELECT * FROM company.test"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d domain.Data
		err := rows.Scan(&d.CompanyID, &d.BranchID, &d.UserID, &d.DataValue, &d.CreatedAt, &d.Username, &d.Password, &d.Role)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *companyRepository) GetMe(data *domain.Data) (*domain.Data, error) {
	query := "SELECT * FROM company.test WHERE company_id = $1 AND branch_id = $2 AND user_id = $3"
	err := r.db.QueryRow(query, data.CompanyID, data.BranchID, data.UserID).Scan(&data.CompanyID, &data.BranchID, &data.UserID, &data.DataValue, &data.CreatedAt, &data.Username, &data.Password, &data.Role)
	if err != nil {
		return nil, err
	}

	return data, nil
}
