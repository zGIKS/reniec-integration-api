package appointments

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAppointment(appointment *Appointment) error {
	query := `
		INSERT INTO appointments (client_id, service_id, appointment_date, start_time, end_time, attended_by, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		appointment.ClientID,
		appointment.ServiceID,
		appointment.AppointmentDate,
		appointment.StartTime,
		appointment.EndTime,
		appointment.AttendedBy,
		appointment.Status,
	).Scan(
		&appointment.ID,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating appointment: %w", err)
	}

	return nil
}

func (r *Repository) GetAppointmentByID(id string) (*Appointment, error) {
	appointment := &Appointment{}
	query := `
		SELECT id, client_id, service_id, appointment_date, start_time, end_time, 
		       attended_by, status, cancelled_by, cancelled_by_type, cancellation_reason,
		       created_at, updated_at
		FROM appointments WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&appointment.ID,
		&appointment.ClientID,
		&appointment.ServiceID,
		&appointment.AppointmentDate,
		&appointment.StartTime,
		&appointment.EndTime,
		&appointment.AttendedBy,
		&appointment.Status,
		&appointment.CancelledBy,
		&appointment.CancelledByType,
		&appointment.CancellationReason,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("appointment not found")
		}
		return nil, fmt.Errorf("error getting appointment: %w", err)
	}

	return appointment, nil
}

func (r *Repository) GetAppointmentWithDetails(id string) (*AppointmentWithDetails, error) {
	appointment := &AppointmentWithDetails{}
	query := `
		SELECT a.id, a.client_id, a.service_id, a.appointment_date, a.start_time, a.end_time, 
		       a.attended_by, a.status, a.cancelled_by, a.cancelled_by_type, a.cancellation_reason,
		       a.created_at, a.updated_at,
		       CONCAT(c.first_name, ' ', c.last_name) as client_name, c.dni as client_dni,
		       s.name as service_name, s.price as service_price, s.duration_minutes as service_duration
		FROM appointments a
		JOIN clients c ON a.client_id = c.id
		JOIN services s ON a.service_id = s.id
		WHERE a.id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&appointment.ID,
		&appointment.ClientID,
		&appointment.ServiceID,
		&appointment.AppointmentDate,
		&appointment.StartTime,
		&appointment.EndTime,
		&appointment.AttendedBy,
		&appointment.Status,
		&appointment.CancelledBy,
		&appointment.CancelledByType,
		&appointment.CancellationReason,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
		&appointment.ClientName,
		&appointment.ClientDNI,
		&appointment.ServiceName,
		&appointment.ServicePrice,
		&appointment.ServiceDuration,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("appointment not found")
		}
		return nil, fmt.Errorf("error getting appointment with details: %w", err)
	}

	return appointment, nil
}

func (r *Repository) CheckAvailability(date time.Time, startTime string, attendedBy string, excludeAppointmentID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM appointments 
		WHERE appointment_date = $1 AND start_time = $2 AND attended_by = $3 AND status != 'cancelled'`
	
	args := []interface{}{date, startTime, attendedBy}
	
	if excludeAppointmentID != "" {
		query += " AND id != $4"
		args = append(args, excludeAppointmentID)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking availability: %w", err)
	}

	return count == 0, nil
}

func (r *Repository) GetAppointmentsByDateRange(startDate, endDate time.Time) ([]AppointmentWithDetails, error) {
	query := `
		SELECT a.id, a.client_id, a.service_id, a.appointment_date, a.start_time, a.end_time, 
		       a.attended_by, a.status, a.cancelled_by, a.cancelled_by_type, a.cancellation_reason,
		       a.created_at, a.updated_at,
		       CONCAT(c.first_name, ' ', c.last_name) as client_name, c.dni as client_dni,
		       s.name as service_name, s.price as service_price, s.duration_minutes as service_duration
		FROM appointments a
		JOIN clients c ON a.client_id = c.id
		JOIN services s ON a.service_id = s.id
		WHERE a.appointment_date BETWEEN $1 AND $2
		ORDER BY a.appointment_date ASC, a.start_time ASC`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying appointments by date range: %w", err)
	}
	defer rows.Close()

	var appointments []AppointmentWithDetails
	for rows.Next() {
		var appointment AppointmentWithDetails
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.ServiceID,
			&appointment.AppointmentDate,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.AttendedBy,
			&appointment.Status,
			&appointment.CancelledBy,
			&appointment.CancelledByType,
			&appointment.CancellationReason,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
			&appointment.ClientName,
			&appointment.ClientDNI,
			&appointment.ServiceName,
			&appointment.ServicePrice,
			&appointment.ServiceDuration,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning appointment: %w", err)
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *Repository) GetAppointmentsByClient(clientID string) ([]AppointmentWithDetails, error) {
	query := `
		SELECT a.id, a.client_id, a.service_id, a.appointment_date, a.start_time, a.end_time, 
		       a.attended_by, a.status, a.cancelled_by, a.cancelled_by_type, a.cancellation_reason,
		       a.created_at, a.updated_at,
		       CONCAT(c.first_name, ' ', c.last_name) as client_name, c.dni as client_dni,
		       s.name as service_name, s.price as service_price, s.duration_minutes as service_duration
		FROM appointments a
		JOIN clients c ON a.client_id = c.id
		JOIN services s ON a.service_id = s.id
		WHERE a.client_id = $1
		ORDER BY a.appointment_date DESC, a.start_time DESC`

	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("error querying appointments by client: %w", err)
	}
	defer rows.Close()

	var appointments []AppointmentWithDetails
	for rows.Next() {
		var appointment AppointmentWithDetails
		err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.ServiceID,
			&appointment.AppointmentDate,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.AttendedBy,
			&appointment.Status,
			&appointment.CancelledBy,
			&appointment.CancelledByType,
			&appointment.CancellationReason,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
			&appointment.ClientName,
			&appointment.ClientDNI,
			&appointment.ServiceName,
			&appointment.ServicePrice,
			&appointment.ServiceDuration,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning appointment: %w", err)
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *Repository) UpdateAppointment(id string, updates UpdateAppointmentRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.AppointmentDate != nil {
		setParts = append(setParts, fmt.Sprintf("appointment_date = $%d", argIndex))
		args = append(args, *updates.AppointmentDate)
		argIndex++
	}
	if updates.StartTime != nil {
		setParts = append(setParts, fmt.Sprintf("start_time = $%d", argIndex))
		args = append(args, *updates.StartTime)
		argIndex++
	}
	if updates.AttendedBy != nil {
		setParts = append(setParts, fmt.Sprintf("attended_by = $%d", argIndex))
		args = append(args, updates.AttendedBy)
		argIndex++
	}
	if updates.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *updates.Status)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE appointments SET %s WHERE id = $%d", 
		strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) CancelAppointment(id string, cancelledBy, cancelledByType, reason string) error {
	query := `
		UPDATE appointments 
		SET status = 'cancelled', cancelled_by = $2, cancelled_by_type = $3, cancellation_reason = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status != 'cancelled'`

	result, err := r.db.Exec(query, id, cancelledBy, cancelledByType, reason)
	if err != nil {
		return fmt.Errorf("error cancelling appointment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("appointment not found or already cancelled")
	}

	return nil
}

