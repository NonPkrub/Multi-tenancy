package repositories

import (
	"errors"
	"fmt"
	"go-multi-tenancy/internals/core/domain"

	"github.com/jmoiron/sqlx"
)

type manageRepository struct {
	db *sqlx.DB
}

func NewManageRepository(db *sqlx.DB) *manageRepository {
	return &manageRepository{db: db}
}

func (m *manageRepository) GetCompany() ([]domain.Manage, error) {
	query := `SELECT inhrelid::regclass AS company
		FROM pg_inherits
		WHERE inhparent = 'company.onesystem'::regclass;`

	var companies []domain.Manage
	err := m.db.Select(&companies, query)
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (m *manageRepository) GetBranch(data *domain.Manage) ([]domain.Manage, error) {

	query := `SELECT inhrelid::regclass AS branch, inhparent::regclass AS company
		FROM pg_inherits
		WHERE inhparent = ($1)::regclass;`

	var branch []domain.Manage
	err := m.db.Get(&branch, query, data.Company)
	if err != nil {
		return nil, err
	}

	return branch, nil

}

func (m *manageRepository) CreateCompany(data *domain.Manage) (*domain.Manage, error) {
	exists, err := m.tableExists(data.Company)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("the company already exist")
	}

	query := fmt.Sprintf(`CREATE TABLE company.%s PARTITION OF company.onesystem
    FOR VALUES IN ('%s')
   	PARTITION BY LIST (branch) ;`, data.Company, data.Company)

	_, err = m.db.Exec(query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *manageRepository) CreateBranch(data *domain.Manage) (*domain.Manage, error) {
	exists, err := m.tableExists(data.Company)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("the company does not exist")
	}

	query := fmt.Sprintf(`CREATE TABLE company.%s PARTITION OF company.%s
    FOR VALUES IN ('%s');`, data.Branch, data.Company, data.Branch)

	_, err = m.db.Exec(query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *manageRepository) DeleteCompany(data *domain.Manage) error {

	exists, err := m.tableExists(data.Company)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("partitioned table does not exist")
	}

	query := "DROP TABLE " + data.Company
	_, err = m.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (m *manageRepository) DeleteBranch(data *domain.Manage) error {
	exists, err := m.tableExists(data.Branch)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("partitioned table does not exist")
	}

	query := "DROP TABLE " + data.Branch
	_, err = m.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (m *manageRepository) tableExists(tableName string) (bool, error) {
	query := `
         SELECT EXISTS (
            SELECT * FROM information_schema.tables
            WHERE table_schema = 'company' 
            AND table_name = $1
        )`
	var exists bool
	err := m.db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
