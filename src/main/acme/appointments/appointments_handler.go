package appointments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppointmentsHandler struct {
	service *AppointmentService
}

func NewAppointmentsHandler(service *AppointmentService) *AppointmentsHandler {
	return &AppointmentsHandler{service: service}
}

// CreateAppointment godoc
// @Summary Create a new appointment
// @Description Create a new appointment with availability validation
// @Tags appointments
// @Accept json
// @Produce json
// @Param appointment body CreateAppointmentRequest true "Appointment data"
// @Success 201 {object} Appointment
// @Failure 400 {object} map[string]interface{}
// @Router /appointments [post]
func (h *AppointmentsHandler) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.service.CreateAppointment(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

func (h *AppointmentsHandler) GetAppointmentByID(c *gin.Context) {
	id := c.Param("id")

	appointment, err := h.service.GetAppointmentByID(id)
	if err != nil {
		if err.Error() == "appointment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *AppointmentsHandler) GetAppointmentWithDetails(c *gin.Context) {
	id := c.Param("id")

	appointment, err := h.service.GetAppointmentWithDetails(id)
	if err != nil {
		if err.Error() == "appointment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *AppointmentsHandler) UpdateAppointment(c *gin.Context) {
	id := c.Param("id")

	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.service.UpdateAppointment(id, req)
	if err != nil {
		if err.Error() == "appointment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *AppointmentsHandler) GetAppointmentsByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date query parameters are required (YYYY-MM-DD format)"})
		return
	}

	appointments, err := h.service.GetAppointmentsByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentsHandler) GetAppointmentsByClient(c *gin.Context) {
	clientID := c.Param("client_id")

	appointments, err := h.service.GetAppointmentsByClient(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *AppointmentsHandler) CheckAvailability(c *gin.Context) {
	date := c.Query("date")
	startTime := c.Query("start_time")
	attendedBy := c.Query("attended_by")

	if date == "" || startTime == "" || attendedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date, start_time, and attended_by query parameters are required"})
		return
	}

	available, err := h.service.CheckAvailability(date, startTime, attendedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"date":      date,
		"start_time": startTime,
		"attended_by": attendedBy,
	})
}

// CancelAppointment godoc
// @Summary Cancel an appointment
// @Description Cancel an appointment with reason and audit trail
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID"
// @Param cancellation body CancelAppointmentRequest true "Cancellation data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /appointments/{id}/cancel [put]
func (h *AppointmentsHandler) CancelAppointment(c *gin.Context) {
	id := c.Param("id")

	var req CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CancelAppointment(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled successfully"})
}

// CancelAppointmentByClient godoc
// @Summary Cancel appointment by client (for chat integration)
// @Description Cancel an appointment from client side (chat/external)
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID"
// @Param cancellation body map[string]string true "Client DNI and reason"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /appointments/{id}/cancel-by-client [put]
func (h *AppointmentsHandler) CancelAppointmentByClient(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ClientDNI string `json:"client_dni" binding:"required"`
		Reason    string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CancelByClient(id, req.ClientDNI, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled by client successfully"})
}

// CancelAppointmentByEmployee godoc
// @Summary Cancel appointment by employee
// @Description Cancel an appointment from employee/backend side
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID"
// @Param cancellation body map[string]string true "Employee email and reason"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /appointments/{id}/cancel-by-employee [put]
func (h *AppointmentsHandler) CancelAppointmentByEmployee(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		EmployeeEmail string `json:"employee_email" binding:"required"`
		Reason        string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CancelByEmployee(id, req.EmployeeEmail, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled by employee successfully"})
}

