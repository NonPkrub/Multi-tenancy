package domain

import (
	"time"

	"github.com/google/uuid"
)

type Data struct {
	Company   string    `json:"company"`
	Branch    string    `json:"branch"`
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
	DeleteAt  time.Time `json:"delete_at"`
	Role      string    `json:"role"`
}

func NewData(company string, branch string, id uuid.UUID, first_name string, last_name string, username string, password string, create_at time.Time, update_at time.Time, delete_at time.Time, role string) *Data {
	return &Data{
		Company:   company,
		Branch:    branch,
		ID:        id,
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

type RegisterInput struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Company   string `json:"company"`
	Branch    string `json:"branch"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Company  string `json:"company"`
	Branch   string `json:"branch"`
}

type DataReply struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Company   string    `json:"company"`
	Branch    string    `json:"branch"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

type DataInput struct {
	Company string `json:"company"`
	Branch  string `json:"branch"`
}

type DataUpdate struct {
	Company   string    `json:"company"`
	Branch    string    `json:"branch"`
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Password  *string   `json:"password"`
}

type DataDelete struct {
	Company string    `json:"company"`
	Branch  string    `json:"branch"`
	ID      uuid.UUID `json:"id"`
}

type Admin struct {
	Company   string `json:"company"`
	Branch    string `json:"branch"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Me struct {
	Company string    `json:"company"`
	Branch  string    `json:"branch"`
	ID      uuid.UUID `json:"id"`
}
