package catalog

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateService(service *Service) error {
	query := `
		INSERT INTO services (name, price, duration_minutes, description, benefits, 
		                     recommended_frequency, includes, contraindications)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		service.Name,
		service.Price,
		service.DurationMinutes,
		service.Description,
		service.Benefits,
		service.RecommendedFrequency,
		service.Includes,
		service.Contraindications,
	).Scan(
		&service.ID,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating service: %w", err)
	}

	return nil
}

func (r *Repository) GetServiceByID(id string) (*Service, error) {
	service := &Service{}
	query := `
		SELECT id, name, price, duration_minutes, description, benefits, 
		       recommended_frequency, includes, contraindications, 
		       created_at, updated_at
		FROM services WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Price,
		&service.DurationMinutes,
		&service.Description,
		&service.Benefits,
		&service.RecommendedFrequency,
		&service.Includes,
		&service.Contraindications,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service not found")
		}
		return nil, fmt.Errorf("error getting service: %w", err)
	}

	return service, nil
}

func (r *Repository) GetAllServices() ([]Service, error) {
	query := `
		SELECT id, name, price, duration_minutes, description, benefits, 
		       recommended_frequency, includes, contraindications, 
		       created_at, updated_at
		FROM services ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying services: %w", err)
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Price,
			&service.DurationMinutes,
			&service.Description,
			&service.Benefits,
			&service.RecommendedFrequency,
			&service.Includes,
			&service.Contraindications,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning service: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *Repository) UpdateService(id string, updates UpdateServiceRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *updates.Name)
		argIndex++
	}
	if updates.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *updates.Price)
		argIndex++
	}
	if updates.DurationMinutes != nil {
		setParts = append(setParts, fmt.Sprintf("duration_minutes = $%d", argIndex))
		args = append(args, *updates.DurationMinutes)
		argIndex++
	}
	if updates.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, updates.Description)
		argIndex++
	}
	if updates.Benefits != nil {
		setParts = append(setParts, fmt.Sprintf("benefits = $%d", argIndex))
		args = append(args, updates.Benefits)
		argIndex++
	}
	if updates.RecommendedFrequency != nil {
		setParts = append(setParts, fmt.Sprintf("recommended_frequency = $%d", argIndex))
		args = append(args, updates.RecommendedFrequency)
		argIndex++
	}
	if updates.Includes != nil {
		setParts = append(setParts, fmt.Sprintf("includes = $%d", argIndex))
		args = append(args, updates.Includes)
		argIndex++
	}
	if updates.Contraindications != nil {
		setParts = append(setParts, fmt.Sprintf("contraindications = $%d", argIndex))
		args = append(args, updates.Contraindications)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE services SET %s WHERE id = $%d", 
		strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) DeleteService(id string) error {
	query := `DELETE FROM services WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("service not found")
	}

	return nil
}

func (r *Repository) GetServicesByPriceRange(minPrice, maxPrice float64) ([]Service, error) {
	query := `
		SELECT id, name, price, duration_minutes, description, benefits, 
		       recommended_frequency, includes, contraindications, 
		       created_at, updated_at
		FROM services 
		WHERE price BETWEEN $1 AND $2 
		ORDER BY price ASC`

	rows, err := r.db.Query(query, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("error querying services by price range: %w", err)
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Price,
			&service.DurationMinutes,
			&service.Description,
			&service.Benefits,
			&service.RecommendedFrequency,
			&service.Includes,
			&service.Contraindications,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning service: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}