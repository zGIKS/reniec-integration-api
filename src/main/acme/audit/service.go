package audit

import (
	"encoding/json"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) LogAction(req CreateAuditLogRequest) error {
	var oldValuesJSON, newValuesJSON json.RawMessage
	var err error

	if req.OldValues != nil {
		oldValuesJSON, err = json.Marshal(req.OldValues)
		if err != nil {
			return fmt.Errorf("error marshaling old values: %w", err)
		}
	}

	if req.NewValues != nil {
		newValuesJSON, err = json.Marshal(req.NewValues)
		if err != nil {
			return fmt.Errorf("error marshaling new values: %w", err)
		}
	}

	log := &AuditLog{
		TableName:     req.TableName,
		RecordID:      req.RecordID,
		Action:        string(req.Action),
		OldValues:     oldValuesJSON,
		NewValues:     newValuesJSON,
		ChangedBy:     req.ChangedBy,
		ChangedByType: string(req.ChangedByType),
		Reason:        req.Reason,
	}

	return s.repo.CreateAuditLog(log)
}

func (s *Service) GetAuditHistory(tableName, recordID string) ([]AuditLog, error) {
	return s.repo.GetAuditLogsByRecordID(tableName, recordID)
}

func (s *Service) GetAuditLogsByDateRange(tableName, startDate, endDate string) ([]AuditLog, error) {
	return s.repo.GetAuditLogsByDateRange(tableName, startDate, endDate)
}