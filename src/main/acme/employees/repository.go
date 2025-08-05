package employees

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

func (r *Repository) GetAllEmployees() ([]Employee, error) {
	query := `
		SELECT id, name, paternal_surname, maternal_surname, role, phone, email, 
		       created_at, updated_at
		FROM employees 
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying employees: %w", err)
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var employee Employee
		err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.PaternalSurname,
			&employee.MaternalSurname,
			&employee.Role,
			&employee.Phone,
			&employee.Email,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee: %w", err)
		}
		employees = append(employees, employee)
	}

	return employees, nil
}


func (r *Repository) GetEmployeeByID(id string) (*Employee, error) {
	employee := &Employee{}
	query := `
		SELECT id, name, paternal_surname, maternal_surname, role, phone, email, 
		       created_at, updated_at
		FROM employees WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&employee.ID,
		&employee.Name,
		&employee.PaternalSurname,
		&employee.MaternalSurname,
		&employee.Role,
		&employee.Phone,
		&employee.Email,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("error getting employee: %w", err)
	}

	return employee, nil
}
