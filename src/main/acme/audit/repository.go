package audit

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

func (r *Repository) CreateAuditLog(log *AuditLog) error {
	query := `
		INSERT INTO audit_logs (table_name, record_id, action, old_values, new_values, 
		                       changed_by, changed_by_type, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		log.TableName,
		log.RecordID,
		log.Action,
		log.OldValues,
		log.NewValues,
		log.ChangedBy,
		log.ChangedByType,
		log.Reason,
	).Scan(
		&log.ID,
		&log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating audit log: %w", err)
	}

	return nil
}

func (r *Repository) GetAuditLogsByRecordID(tableName, recordID string) ([]AuditLog, error) {
	query := `
		SELECT id, table_name, record_id, action, old_values, new_values, 
		       changed_by, changed_by_type, reason, created_at
		FROM audit_logs 
		WHERE table_name = $1 AND record_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, tableName, recordID)
	if err != nil {
		return nil, fmt.Errorf("error querying audit logs: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.TableName,
			&log.RecordID,
			&log.Action,
			&log.OldValues,
			&log.NewValues,
			&log.ChangedBy,
			&log.ChangedByType,
			&log.Reason,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *Repository) GetAuditLogsByDateRange(tableName string, startDate, endDate string) ([]AuditLog, error) {
	query := `
		SELECT id, table_name, record_id, action, old_values, new_values, 
		       changed_by, changed_by_type, reason, created_at
		FROM audit_logs 
		WHERE table_name = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, tableName, startDate+" 00:00:00", endDate+" 23:59:59")
	if err != nil {
		return nil, fmt.Errorf("error querying audit logs by date range: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.TableName,
			&log.RecordID,
			&log.Action,
			&log.OldValues,
			&log.NewValues,
			&log.ChangedBy,
			&log.ChangedByType,
			&log.Reason,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
