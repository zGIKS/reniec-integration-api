package iam

import (
	"time"
)

type Client struct {
	ID               string    `json:"id" db:"id"`
	FirstName        string    `json:"first_name" db:"first_name"`
	LastName         string    `json:"last_name" db:"last_name"`
	SecondLastName   *string   `json:"second_last_name" db:"second_last_name"`
	DNI              string    `json:"dni" db:"dni"`
	Email            string    `json:"email" db:"email"`
	Phone            *string   `json:"phone" db:"phone"`
	RegistrationDate time.Time `json:"registration_date" db:"registration_date"`
	ReniecValidated  bool      `json:"reniec_validated" db:"reniec_validated"`
	FullName         string    `json:"full_name"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type CreateClientRequest struct {
	FirstName      string  `json:"first_name" binding:"required"`
	LastName       string  `json:"last_name" binding:"required"`
	SecondLastName *string `json:"second_last_name"`
	DNI            string  `json:"dni" binding:"required,len=8"`
	Email          string  `json:"email" binding:"required,email"`
	Phone          *string `json:"phone"`
}

type UpdateClientRequest struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	SecondLastName *string `json:"second_last_name"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
}

type ReniecResponse struct {
	FirstName        string `json:"first_name"`
	FirstLastName    string `json:"first_last_name"`
	SecondLastName   string `json:"second_last_name"`
	FullName         string `json:"full_name"`
	DocumentNumber   string `json:"document_number"`
}

type ReniecValidationResult struct {
	IsValid bool           `json:"is_valid"`
	Data    ReniecResponse `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

type ReniecValidationRequest struct {
	DNI string `json:"dni" binding:"required,len=8"`
}

func (c *Client) GenerateFullName() {
	fullName := c.FirstName + " " + c.LastName
	if c.SecondLastName != nil && *c.SecondLastName != "" {
		fullName += " " + *c.SecondLastName
	}
	c.FullName = fullName
}