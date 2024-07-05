package repositories

import (
	"go-multi-tenancy/internals/core/domain"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type companyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *companyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Register(data *domain.Data) (*domain.Data, error) {
	query := "INSERT INTO company.onesystem (company, branch, first_name, last_name, username, password, role)   VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING company, branch,  id, first_name, last_name, username,  create_at"
	err := r.db.QueryRow(query, data.Company, data.Branch, data.FirstName, data.LastName, data.Username, data.Password, data.Role).Scan(&data.Company, &data.Branch, &data.ID, &data.FirstName, &data.LastName, &data.Username, &data.CreateAt)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) Login(data *domain.Data) (*domain.Data, error) {
	query := "SELECT * FROM company.onesystem WHERE username = $1 AND company = $2 AND branch = $3"
	err := r.db.QueryRow(query, data.Username, data.Company, data.Branch).Scan(&data.Company, &data.Branch, &data.ID, &data.FirstName, &data.LastName, &data.Username, &data.Password, &data.CreateAt, &data.UpdateAt, &data.DeleteAt, &data.Role)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) GetData(data *domain.Data) (*domain.Data, error) {
	query := " SELECT * FROM company.onesystem WHERE company = $1 AND branch = $2 "
	err := r.db.QueryRow(query, data.Company, data.Branch).Scan(&data.Company, &data.Branch, &data.ID, &data.FirstName, &data.LastName, &data.Username, &data.Password, &data.CreateAt, &data.UpdateAt, &data.DeleteAt, &data.Role)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) UpdateData(data *domain.Data) (*domain.Data, error) {

	var fields []string
	var args []interface{}
	argIndex := 1

	if data.FirstName != "" {
		fields = append(fields, "first_name = $"+strconv.Itoa(argIndex))
		args = append(args, data.FirstName)
		argIndex++
	}
	if data.LastName != "" {
		fields = append(fields, "last_name = $"+strconv.Itoa(argIndex))
		args = append(args, data.LastName)
		argIndex++
	}
	if data.Password != "" {
		fields = append(fields, "password = $"+strconv.Itoa(argIndex))
		args = append(args, data.Password)
		argIndex++
	}

	query := "UPDATE company.onesystem SET " + strings.Join(fields, ", ") + " WHERE company = $" + strconv.Itoa(argIndex) + " AND branch = $" + strconv.Itoa(argIndex+1) + " AND username = $" + strconv.Itoa(argIndex+2)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	updateData, err := r.GetOne(data)
	if err != nil {
		return nil, err
	}

	return updateData, nil

}

func (r *companyRepository) GetOne(data *domain.Data) (*domain.Data, error) {
	query := "SELECT * FROM company.onesystem WHERE company = $1 AND branch = $2 AND id = $3"
	err := r.db.QueryRow(query, data.Company, data.Branch, data.ID).Scan(&data.Company, &data.Branch, &data.ID, &data.FirstName, &data.LastName, &data.Username, &data.Password, &data.CreateAt, &data.UpdateAt, &data.DeleteAt, &data.Role)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *companyRepository) DeleteData(data *domain.Data) error {
	query := "DELETE FROM company.onesystem WHERE company = $1 AND branch = $2 AND id = $3"
	_, err := r.db.Exec(query, data.Company, data.Branch, data.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *companyRepository) GetAllData() ([]domain.Data, error) {
	data := []domain.Data{}

	query := "SELECT * FROM company.onesystem"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d domain.Data
		err := rows.Scan(&d.Company, &d.Branch, &d.ID, &d.FirstName, &d.LastName, &d.Username, &d.Password, &d.CreateAt, &d.UpdateAt, &d.DeleteAt, &d.Role)
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
	query := "SELECT * FROM company.onesystem WHERE company = $1 AND branch = $2 AND id = $3"
	err := r.db.QueryRow(query, data.Company, data.Branch, data.ID).Scan(&data.Company, &data.Branch, &data.ID, &data.FirstName, &data.LastName, &data.Username, &data.Password, &data.CreateAt, &data.UpdateAt, &data.DeleteAt, &data.Role)
	if err != nil {
		return nil, err
	}

	return data, nil
}
