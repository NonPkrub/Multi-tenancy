package domain

import (
	"time"
)

type Data struct {
	CompanyID int       `json:"company_id"`
	BranchID  int       `json:"branch_id"`
	UserID    int       `json:"user_id"`
	DataValue string    `json:"data_value"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
}

func NewData(company_id int, branch_id int, user_id int, data_value string, created_at time.Time, username string, password string, role string) *Data {
	return &Data{
		CompanyID: company_id,
		BranchID:  branch_id,
		UserID:    user_id,
		DataValue: data_value,
		CreatedAt: created_at,
		Username:  username,
		Password:  password,
		Role:      role,
	}
}

type RegisterInput struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
	DataValue string `json:"data_value"`
}

type LoginInput struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
}

type DataReply struct {
	Username  string    `json:"username"`
	UserID    int       `json:"user_id"`
	DataValue string    `json:"data_value"`
	CreatedAt time.Time `json:"created_at"`
	CompanyID int       `json:"company_id"`
	BranchID  int       `json:"branch_id"`
}

type DataInput struct {
	CompanyID int `json:"company_id"`
	BranchID  int `json:"branch_id"`
}

type DataUpdate struct {
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
	UserID    int    `json:"user_id"`
	DataValue string `json:"data_value"`
}

type DataDelete struct {
	CompanyID int `json:"company_id"`
	BranchID  int `json:"branch_id"`
	UserID    int `json:"user_id"`
}

type Admin struct {
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
	UserID    int    `json:"user_id"`
	DataValue string `json:"data_value"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Me struct {
	CompanyID int `json:"company_id"`
	BranchID  int `json:"branch_id"`
	UserID    int `json:"user_id"`
}
