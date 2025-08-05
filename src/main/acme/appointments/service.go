package appointments

import (
	"acme/audit"
	"fmt"
	"strconv"
	"time"
)

type AppointmentService struct {
	repo        *Repository
	auditService *audit.Service
}

func NewService(repo *Repository, auditService *audit.Service) *AppointmentService {
	return &AppointmentService{
		repo:        repo,
		auditService: auditService,
	}
}

func (s *AppointmentService) CreateAppointment(req CreateAppointmentRequest) (*Appointment, error) {
	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid appointment date format, use YYYY-MM-DD: %w", err)
	}

	if appointmentDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, fmt.Errorf("cannot create appointments for past dates")
	}

	if !s.isValidTimeFormat(req.StartTime) {
		return nil, fmt.Errorf("invalid start time format, use HH:MM")
	}

	available, err := s.repo.CheckAvailability(appointmentDate, req.StartTime, req.AttendedBy, "")
	if err != nil {
		return nil, fmt.Errorf("error checking availability: %w", err)
	}

	if !available {
		return nil, fmt.Errorf("the requested time slot is not available")
	}

	endTime, err := s.calculateEndTime(req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("error calculating end time: %w", err)
	}

	appointment := &Appointment{
		ClientID:        req.ClientID,
		ServiceID:       req.ServiceID,
		AppointmentDate: appointmentDate,
		StartTime:       req.StartTime,
		EndTime:         endTime,
		AttendedBy:      &req.AttendedBy,
		Status:          string(StatusPending),
	}

	if err := s.repo.CreateAppointment(appointment); err != nil {
		return nil, fmt.Errorf("error creating appointment: %w", err)
	}

	return appointment, nil
}

func (s *AppointmentService) calculateEndTime(startTime string) (string, error) {
	startHour, startMin, err := s.parseTime(startTime)
	if err != nil {
		return "", err
	}

	totalMinutes := startHour*60 + startMin + 60

	endHour := totalMinutes / 60
	endMin := totalMinutes % 60

	if endHour >= 24 {
		endHour = endHour % 24
	}

	return fmt.Sprintf("%02d:%02d", endHour, endMin), nil
}

func (s *AppointmentService) parseTime(timeStr string) (int, int, error) {
	parts := make([]string, 2)
	if len(timeStr) == 5 && timeStr[2] == ':' {
		parts[0] = timeStr[:2]
		parts[1] = timeStr[3:]
	} else {
		return 0, 0, fmt.Errorf("invalid time format")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("invalid hour")
	}

	min, err := strconv.Atoi(parts[1])
	if err != nil || min < 0 || min > 59 {
		return 0, 0, fmt.Errorf("invalid minute")
	}

	return hour, min, nil
}

func (s *AppointmentService) isValidTimeFormat(timeStr string) bool {
	_, _, err := s.parseTime(timeStr)
	return err == nil
}

func (s *AppointmentService) GetAppointmentByID(id string) (*Appointment, error) {
	return s.repo.GetAppointmentByID(id)
}

func (s *AppointmentService) GetAppointmentWithDetails(id string) (*AppointmentWithDetails, error) {
	return s.repo.GetAppointmentWithDetails(id)
}

func (s *AppointmentService) UpdateAppointment(id string, req UpdateAppointmentRequest) (*Appointment, error) {
	if req.Status != nil {
		status := AppointmentStatus(*req.Status)
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid status: %s", *req.Status)
		}
	}

	if req.AppointmentDate != nil {
		appointmentDate, err := time.Parse("2006-01-02", *req.AppointmentDate)
		if err != nil {
			return nil, fmt.Errorf("invalid appointment date format, use YYYY-MM-DD: %w", err)
		}

		if appointmentDate.Before(time.Now().Truncate(24 * time.Hour)) {
			return nil, fmt.Errorf("cannot update to past dates")
		}
	}

	if req.StartTime != nil && !s.isValidTimeFormat(*req.StartTime) {
		return nil, fmt.Errorf("invalid start time format, use HH:MM")
	}

	if (req.AppointmentDate != nil || req.StartTime != nil || req.AttendedBy != nil) && req.Status == nil {
		currentAppointment, err := s.repo.GetAppointmentByID(id)
		if err != nil {
			return nil, err
		}

		date := currentAppointment.AppointmentDate
		if req.AppointmentDate != nil {
			date, _ = time.Parse("2006-01-02", *req.AppointmentDate)
		}

		startTime := currentAppointment.StartTime
		if req.StartTime != nil {
			startTime = *req.StartTime
		}

		attendedBy := ""
		if currentAppointment.AttendedBy != nil {
			attendedBy = *currentAppointment.AttendedBy
		}
		if req.AttendedBy != nil {
			attendedBy = *req.AttendedBy
		}

		available, err := s.repo.CheckAvailability(date, startTime, attendedBy, id)
		if err != nil {
			return nil, fmt.Errorf("error checking availability: %w", err)
		}

		if !available {
			return nil, fmt.Errorf("the requested time slot is not available")
		}
	}

	if err := s.repo.UpdateAppointment(id, req); err != nil {
		return nil, fmt.Errorf("error updating appointment: %w", err)
	}

	return s.repo.GetAppointmentByID(id)
}

func (s *AppointmentService) GetAppointmentsByDateRange(startDate, endDate string) ([]AppointmentWithDetails, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format, use YYYY-MM-DD: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format, use YYYY-MM-DD: %w", err)
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	return s.repo.GetAppointmentsByDateRange(start, end)
}

func (s *AppointmentService) GetAppointmentsByClient(clientID string) ([]AppointmentWithDetails, error) {
	return s.repo.GetAppointmentsByClient(clientID)
}


func (s *AppointmentService) CheckAvailability(date, startTime, attendedBy string) (bool, error) {
	appointmentDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false, fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
	}

	if !s.isValidTimeFormat(startTime) {
		return false, fmt.Errorf("invalid start time format, use HH:MM")
	}

	return s.repo.CheckAvailability(appointmentDate, startTime, attendedBy, "")
}

func (s *AppointmentService) CancelAppointment(id string, req CancelAppointmentRequest) error {
	cancelledByType := CancelledByType(req.CancelledByType)
	if !cancelledByType.IsValid() {
		return fmt.Errorf("invalid cancelled_by_type: %s", req.CancelledByType)
	}

	// Get current appointment for audit log
	currentAppointment, err := s.repo.GetAppointmentByID(id)
	if err != nil {
		return fmt.Errorf("error getting appointment: %w", err)
	}

	if currentAppointment.Status == string(StatusCancelled) {
		return fmt.Errorf("appointment is already cancelled")
	}

	// Cancel the appointment
	err = s.repo.CancelAppointment(id, req.CancelledBy, req.CancelledByType, req.Reason)
	if err != nil {
		return fmt.Errorf("error cancelling appointment: %w", err)
	}

	// Log the cancellation in audit
	auditReq := audit.CreateAuditLogRequest{
		TableName:     "appointments",
		RecordID:      id,
		Action:        audit.ActionCancel,
		OldValues:     currentAppointment,
		NewValues:     map[string]interface{}{
			"status":              "cancelled",
			"cancelled_by":        req.CancelledBy,
			"cancelled_by_type":   req.CancelledByType,
			"cancellation_reason": req.Reason,
		},
		ChangedBy:     req.CancelledBy,
		ChangedByType: audit.ChangedByType(req.CancelledByType),
		Reason:        &req.Reason,
	}

	if err := s.auditService.LogAction(auditReq); err != nil {
		// Log error but don't fail the cancellation
		fmt.Printf("Warning: Failed to log audit entry for appointment cancellation: %v\n", err)
	}

	return nil
}

func (s *AppointmentService) CancelByClient(appointmentID, clientDNI, reason string) error {
	req := CancelAppointmentRequest{
		CancelledBy:     clientDNI,
		CancelledByType: string(CancelledByClient),
		Reason:          reason,
	}
	return s.CancelAppointment(appointmentID, req)
}

func (s *AppointmentService) CancelByEmployee(appointmentID, employeeEmail, reason string) error {
	req := CancelAppointmentRequest{
		CancelledBy:     employeeEmail,
		CancelledByType: string(CancelledByEmployee),
		Reason:          reason,
	}
	return s.CancelAppointment(appointmentID, req)
}