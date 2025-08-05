package audit

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID            string          `json:"id" db:"id"`
	TableName     string          `json:"table_name" db:"table_name"`
	RecordID      string          `json:"record_id" db:"record_id"`
	Action        string          `json:"action" db:"action"`
	OldValues     json.RawMessage `json:"old_values" db:"old_values"`
	NewValues     json.RawMessage `json:"new_values" db:"new_values"`
	ChangedBy     string          `json:"changed_by" db:"changed_by"`
	ChangedByType string          `json:"changed_by_type" db:"changed_by_type"`
	Reason        *string         `json:"reason" db:"reason"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

type AuditAction string

const (
	ActionCreate AuditAction = "CREATE"
	ActionUpdate AuditAction = "UPDATE"
	ActionDelete AuditAction = "DELETE"
	ActionCancel AuditAction = "CANCEL"
)

type ChangedByType string

const (
	ChangedByClient   ChangedByType = "client"
	ChangedByEmployee ChangedByType = "employee"
	ChangedBySystem   ChangedByType = "system"
)

type CreateAuditLogRequest struct {
	TableName     string          `json:"table_name"`
	RecordID      string          `json:"record_id"`
	Action        AuditAction     `json:"action"`
	OldValues     interface{}     `json:"old_values,omitempty"`
	NewValues     interface{}     `json:"new_values,omitempty"`
	ChangedBy     string          `json:"changed_by"`
	ChangedByType ChangedByType   `json:"changed_by_type"`
	Reason        *string         `json:"reason,omitempty"`
}