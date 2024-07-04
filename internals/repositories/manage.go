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

func (m *manageRepository) GetBranch(data *domain.GetBranch) ([]domain.GetBranch, error) {

	query := `SELECT  inhparent::regclass AS company,inhrelid::regclass AS branch
		FROM pg_inherits
		WHERE inhparent = $1::regclass;`

	rows, err := m.db.Query(query, "company."+data.Company)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branchMap := make(map[string][]*string)
	for rows.Next() {
		var branch, company string
		if err := rows.Scan(&branch, &company); err != nil {
			return nil, err
		}
		branchMap[branch] = append(branchMap[branch], &company)
	}

	var branches []domain.GetBranch
	for company, branch := range branchMap {
		branches = append(branches, domain.GetBranch{
			Company: company,
			Branch:  branch,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return branches, nil

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
		return errors.New("the company does not exist")
	}

	query := "DROP TABLE company." + data.Company
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
		return errors.New("the company does not exist")
	}

	query := "DROP TABLE company." + data.Branch
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

func (m *manageRepository) UpdateCompanyToBranch(data *domain.CompanyAndBranch) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	//step 1: detach partition => old company, old branch
	oldBranch := fmt.Sprintf(`ALTER TABLE company.%s DETACH PARTITION company.%s;`, data.OldCompany, data.OldBranch)
	_, err = tx.Exec(oldBranch)
	if err != nil {
		return err
	}

	//step 2: create branch => new branch, new company , branch name
	newBranch := fmt.Sprintf(`CREATE TABLE company.%s PARTITION OF company.%s  FOR VALUES IN ('%s');`, data.NewBranch, data.NewCompany, data.BranchName)
	_, err = tx.Exec(newBranch)
	if err != nil {
		return err
	}

	//step 3: insert data =>  new branch,new company,branch name, old branch
	query := fmt.Sprintf(`INSERT INTO company.%s (company,branch,first_name,last_name,username,password, create_at, update_at,delete_at, role) SELECT '%s','%s',first_name,last_name,username,password, create_at, update_at,delete_at, role FROM company.%s;`, data.NewBranch, data.NewCompany, data.BranchName, data.OldBranch)
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	//step 4: update company => new company, new branch, old company, old branch
	updateQuery := fmt.Sprintf(`UPDATE company.onesystem SET company = '%s', branch = '%s' WHERE company = '%s' AND branch = '%s';`, data.NewCompany, data.NewBranch, data.OldCompany, data.OldBranch)
	_, err = tx.Exec(updateQuery)
	if err != nil {
		return err
	}

	//step 5: delete company => old branch, old company
	deleteBranch := fmt.Sprintf(`DROP TABLE company.%s;`, data.OldBranch)
	_, err = tx.Exec(deleteBranch)
	if err != nil {
		return err
	}

	deleteCompany := fmt.Sprintf(`DROP TABLE company.%s;`, data.OldCompany)
	_, err = tx.Exec(deleteCompany)
	if err != nil {
		return err
	}

	return nil
}

func (m *manageRepository) UpdateBranchToCompany(data *domain.CompanyAndBranch) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// step 1: detach partition => old company,old branch
	oldBranch := fmt.Sprintf(`ALTER TABLE company.%s DETACH PARTITION company.%s;`, data.OldCompany, data.OldBranch)
	_, err = tx.Exec(oldBranch)
	if err != nil {
		return err
	}

	//step 2: create company => new company, new branch
	newCompany := fmt.Sprintf(`CREATE TABLE company.%s PARTITION OF company.onesystem  FOR VALUES IN ('%s') PARTITION BY LIST (branch);`, data.NewCompany, data.NewBranch)
	_, err = tx.Exec(newCompany)
	if err != nil {
		return err
	}

	//step 3: create branch => new branch, new company , new branch name
	initBranch := fmt.Sprintf(`CREATE TABLE company.%s PARTITION OF company.%s FOR VALUES IN ('%s');`, data.NewBranch, data.NewCompany, data.BranchName)
	_, err = tx.Exec(initBranch)
	if err != nil {
		return err
	}

	//step 4: insert data into new partition => new branch , new company, new branch name , old branch
	insertQuery := `
    INSERT INTO company.%s (company, branch, first_name, last_name, username, password, create_at, update_at, delete_at, role) 
    SELECT $1, $2, first_name, last_name, username, password, create_at, update_at, delete_at, role 
    FROM company.%s 
    ON CONFLICT (company, branch) DO UPDATE
    SET
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        username = EXCLUDED.username,
        password = EXCLUDED.password,
        create_at = EXCLUDED.create_at, 	
        update_at = EXCLUDED.update_at,
        delete_at = EXCLUDED.delete_at,
        role = EXCLUDED.role;
	`

	insertQuery = fmt.Sprintf(insertQuery, data.NewBranch, data.OldBranch)

	_, err = tx.Exec(insertQuery, data.NewCompany, data.BranchName)
	if err != nil {
		return err
	}

	//step 5: update data => new company , old company , old branch
	updateQuery := fmt.Sprintf(`UPDATE company.onesystem SET company = '%s' WHERE company = '%s' AND branch = '%s'`, data.NewCompany, data.OldCompany, data.OldBranch)
	_, err = tx.Exec(updateQuery)
	if err != nil {
		return err
	}

	//step 6: delete branch => old branch
	deleteQuery := fmt.Sprintf(`DROP TABLE company.%s`, data.OldBranch)
	_, err = tx.Exec(deleteQuery)
	if err != nil {
		return err
	}

	return nil
}

func (m *manageRepository) UpdateCompanyName(data *domain.RenameCompany) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	//step 1: rename company => old company, new company
	renameQuery := fmt.Sprintf(`ALTER TABLE company.%s RENAME TO %s ;`, data.OldCompany, data.NewCompany)
	_, err = tx.Exec(renameQuery)
	if err != nil {
		return err
	}

	//step 2: detach company => new company
	detachQuery := fmt.Sprintf(`ALTER TABLE company.onesystem DETACH PARTITION company.%s;`, data.NewCompany)
	_, err = tx.Exec(detachQuery)
	if err != nil {
		return err
	}

	//step 3: update company => new company, old company
	updateQuery := fmt.Sprintf(`UPDATE company.%s SET company = '%s' WHERE company = '%s';`, data.NewCompany, data.NewCompany, data.OldCompany)
	_, err = tx.Exec(updateQuery)
	if err != nil {
		return err
	}

	//step 4: attach company => new company
	attachQuery := fmt.Sprintf(`ALTER TABLE company.onesystem ATTACH PARTITION company.%s FOR VALUES IN ('%s');`, data.NewCompany, data.NewCompany)
	_, err = tx.Exec(attachQuery)
	if err != nil {
		return err
	}

	return nil
}

func (m *manageRepository) UpdateBranchName(data *domain.RenameBranch) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	//step 1: rename branch => old branch, new branch
	renameQuery := fmt.Sprintf(`ALTER TABLE company.%s RENAME TO %s;`, data.OldBranch, data.NewBranch)
	_, err = tx.Exec(renameQuery)
	if err != nil {
		return err
	}

	//step 2: detach branch => new branch, company
	detachQuery := fmt.Sprintf(`ALTER TABLE company.%s DETACH PARTITION company.%s;`, data.Company, data.NewBranch)
	_, err = tx.Exec(detachQuery)
	if err != nil {
		return err
	}

	//step 3: update branch => new branch,company, old branch
	updateQuery := fmt.Sprintf(`UPDATE company.%s SET branch = '%s' WHERE company = '%s' AND branch = '%s';`, data.Company, data.NewBranch, data.Company, data.OldBranch)
	_, err = tx.Exec(updateQuery)
	if err != nil {
		return err
	}

	//step 4: attach branch => new branch, company
	attachQuery := fmt.Sprintf(`ALTER TABLE company.%s ATTACH PARTITION company.%s FOR VALUES IN ('%s');`, data.Company, data.NewBranch, data.NewBranch)
	_, err = tx.Exec(attachQuery)
	if err != nil {
		return err
	}

	return nil
}
