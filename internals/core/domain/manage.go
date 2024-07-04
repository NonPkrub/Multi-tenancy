package domain

import "time"

type Manage struct {
	Company   string    `json:"company"`
	Branch    string    `json:"branch"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
	DeleteAt  time.Time `json:"delete_at"`
	Role      string    `json:"role"`
}

func NewManage(company string, branch string, first_name string, last_name string, username string, password string, create_at time.Time, update_at time.Time, delete_at time.Time, role string) *Manage {
	return &Manage{
		Company:   company,
		Branch:    branch,
		FirstName: first_name,
		LastName:  last_name,
		Username:  username,
		Password:  password,
		CreateAt:  create_at,
		UpdateAt:  update_at,
		DeleteAt:  delete_at,
		Role:      role,
	}
}

type CompanyAndBranch struct {
	OldCompany string `json:"company"`
	OldBranch  string `json:"branch"`
	NewCompany string `json:"new_company"`
	NewBranch  string `json:"new_branch"`
	BranchName string `json:"branch_name"`
}

type RenameCompany struct {
	OldCompany string `json:"old_company"`
	NewCompany string `json:"new_company"`
}

type RenameBranch struct {
	Company   string `json:"company"`
	OldBranch string `json:"old_branch"`
	NewBranch string `json:"new_branch"`
}

type GetBranch struct {
	Company string    `json:"company"`
	Branch  []*string `json:"branch"`
}

type CompanyRequest struct {
	Company string `json:"company"`
}

type BranchRequest struct {
	Company string `json:"company"`
	Branch  string `json:"branch"`
}

type Response struct {
	Company string  `json:"company"`
	Branch  *string `json:"branch"`
}

type ResponseCompany struct {
	Company string `json:"company"`
}

type ResponseBranch struct {
	Company string         `json:"company"`
	Branch  []BranchObject `json:"branch"`
}

type BranchObject struct {
	Name string `json:"name"`
}

type CompanyUpdate struct {
	Company   string  `json:"company"`
	Branch    string  `json:"branch"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Password  *string `json:"password"`
}

type BranchUpdate struct {
	Company   string  `json:"company"`
	Branch    string  `json:"branch"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Password  *string `json:"password"`
}
