package employees

import (
	"time"
)

type Employee struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	PaternalSurname  string    `json:"paternal_surname" db:"paternal_surname"`
	MaternalSurname  *string   `json:"maternal_surname" db:"maternal_surname"`
	Role             string    `json:"role" db:"role"`
	Phone            *string   `json:"phone" db:"phone"`
	Email            *string   `json:"email" db:"email"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func (e *Employee) FullName() string {
	fullName := e.Name + " " + e.PaternalSurname
	if e.MaternalSurname != nil && *e.MaternalSurname != "" {
		fullName += " " + *e.MaternalSurname
	}
	return fullName
}

type CreateEmployeeRequest struct {
	Name             string  `json:"name" binding:"required"`
	PaternalSurname  string  `json:"paternal_surname" binding:"required"`
	MaternalSurname  *string `json:"maternal_surname"`
	Role             string  `json:"role"`
	Phone            *string `json:"phone"`
	Email            *string `json:"email"`
}

type UpdateEmployeeRequest struct {
	Name             *string `json:"name"`
	PaternalSurname  *string `json:"paternal_surname"`
	MaternalSurname  *string `json:"maternal_surname"`
	Role             *string `json:"role"`
	Phone            *string `json:"phone"`
	Email            *string `json:"email"`
}