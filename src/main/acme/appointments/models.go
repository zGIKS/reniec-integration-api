package appointments

import (
	"time"
)

type Appointment struct {
	ID                  string    `json:"id" db:"id"`
	ClientID            string    `json:"client_id" db:"client_id"`
	ServiceID           string    `json:"service_id" db:"service_id"`
	AppointmentDate     time.Time `json:"appointment_date" db:"appointment_date"`
	StartTime           string    `json:"start_time" db:"start_time"`
	EndTime             string    `json:"end_time" db:"end_time"`
	AttendedBy          *string   `json:"attended_by" db:"attended_by"`
	Status              string    `json:"status" db:"status"`
	CancelledBy         *string   `json:"cancelled_by" db:"cancelled_by"`
	CancelledByType     *string   `json:"cancelled_by_type" db:"cancelled_by_type"`
	CancellationReason  *string   `json:"cancellation_reason" db:"cancellation_reason"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type AppointmentWithDetails struct {
	Appointment
	ClientName  string  `json:"client_name"`
	ClientDNI   string  `json:"client_dni"`
	ServiceName string  `json:"service_name"`
	ServicePrice float64 `json:"service_price"`
	ServiceDuration int  `json:"service_duration"`
}

type CreateAppointmentRequest struct {
	ClientID        string `json:"client_id" binding:"required"`
	ServiceID       string `json:"service_id" binding:"required"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	StartTime       string `json:"start_time" binding:"required"`
	AttendedBy      string `json:"attended_by"`
}

type UpdateAppointmentRequest struct {
	AppointmentDate *string `json:"appointment_date"`
	StartTime       *string `json:"start_time"`
	AttendedBy      *string `json:"attended_by"`
	Status          *string `json:"status"`
}

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "pending"
	StatusConfirmed AppointmentStatus = "confirmed"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusCompleted AppointmentStatus = "completed"
)

func (s AppointmentStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusCancelled, StatusCompleted:
		return true
	}
	return false
}

type AvailabilitySlot struct {
	Date      time.Time `json:"date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Available bool      `json:"available"`
}

type CancelAppointmentRequest struct {
	CancelledBy     string `json:"cancelled_by" binding:"required"`
	CancelledByType string `json:"cancelled_by_type" binding:"required"`
	Reason          string `json:"reason" binding:"required"`
}

type CancelledByType string

const (
	CancelledByClient   CancelledByType = "client"
	CancelledByEmployee CancelledByType = "employee"
)

func (c CancelledByType) IsValid() bool {
	switch c {
	case CancelledByClient, CancelledByEmployee:
		return true
	}
	return false
}