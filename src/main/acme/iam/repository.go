package iam

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateClient(client *Client) error {
	query := `
		INSERT INTO clients (first_name, last_name, second_last_name, dni, email, phone, reniec_validated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, registration_date, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		client.FirstName,
		client.LastName,
		client.SecondLastName,
		client.DNI,
		client.Email,
		client.Phone,
		client.ReniecValidated,
	).Scan(
		&client.ID,
		&client.RegistrationDate,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == nil {
		client.GenerateFullName()
	}

	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	return nil
}

func (r *Repository) GetClientByID(id string) (*Client, error) {
	client := &Client{}
	query := `
		SELECT id, first_name, last_name, second_last_name, dni, email, phone, 
		       registration_date, reniec_validated, created_at, updated_at
		FROM clients WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&client.ID,
		&client.FirstName,
		&client.LastName,
		&client.SecondLastName,
		&client.DNI,
		&client.Email,
		&client.Phone,
		&client.RegistrationDate,
		&client.ReniecValidated,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == nil {
		client.GenerateFullName()
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	return client, nil
}

func (r *Repository) GetClientByDNI(dni string) (*Client, error) {
	client := &Client{}
	query := `
		SELECT id, first_name, last_name, second_last_name, dni, email, phone, 
		       registration_date, reniec_validated, created_at, updated_at
		FROM clients WHERE dni = $1`

	err := r.db.QueryRow(query, dni).Scan(
		&client.ID,
		&client.FirstName,
		&client.LastName,
		&client.SecondLastName,
		&client.DNI,
		&client.Email,
		&client.Phone,
		&client.RegistrationDate,
		&client.ReniecValidated,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == nil {
		client.GenerateFullName()
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	return client, nil
}

func (r *Repository) UpdateClient(id string, updates UpdateClientRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.FirstName != nil {
		setParts = append(setParts, fmt.Sprintf("first_name = $%d", argIndex))
		args = append(args, *updates.FirstName)
		argIndex++
	}
	if updates.LastName != nil {
		setParts = append(setParts, fmt.Sprintf("last_name = $%d", argIndex))
		args = append(args, *updates.LastName)
		argIndex++
	}
	if updates.SecondLastName != nil {
		setParts = append(setParts, fmt.Sprintf("second_last_name = $%d", argIndex))
		args = append(args, updates.SecondLastName)
		argIndex++
	}
	if updates.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *updates.Email)
		argIndex++
	}
	if updates.Phone != nil {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, updates.Phone)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE clients SET %s WHERE id = $%d", 
		fmt.Sprintf("%s", setParts[0]), argIndex)
	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf("UPDATE clients SET %s, %s WHERE id = $%d", 
			setParts[0], setParts[i], argIndex)
	}
	
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) UpdateReniecValidation(id string, validated bool) error {
	query := `UPDATE clients SET reniec_validated = $1 WHERE id = $2`
	_, err := r.db.Exec(query, validated, id)
	return err
}

func (r *Repository) GetAllClients() ([]Client, error) {
	query := `
		SELECT id, first_name, last_name, second_last_name, dni, email, phone, 
		       registration_date, reniec_validated, created_at, updated_at
		FROM clients ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying clients: %w", err)
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		err := rows.Scan(
			&client.ID,
			&client.FirstName,
			&client.LastName,
			&client.SecondLastName,
			&client.DNI,
			&client.Email,
			&client.Phone,
			&client.RegistrationDate,
			&client.ReniecValidated,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning client: %w", err)
		}
		client.GenerateFullName()
		clients = append(clients, client)
	}

	return clients, nil
}