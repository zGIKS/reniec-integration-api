package iam

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAMHandler struct {
	service *IAMService
}

func NewIAMHandler(service *IAMService) *IAMHandler {
	return &IAMHandler{service: service}
}

// CreateClient godoc
// @Summary Create a new client
// @Description Create a new client with RENIEC validation
// @Tags clients
// @Accept json
// @Produce json
// @Param client body CreateClientRequest true "Client data"
// @Success 201 {object} Client
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /clients [post]
func (h *IAMHandler) CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.service.CreateClient(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, client)
}

func (h *IAMHandler) GetClientByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	client, err := h.service.GetClientByID(id)
	if err != nil {
		if err.Error() == "client not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *IAMHandler) GetClientByDNI(c *gin.Context) {
	dni := c.Param("dni")
	if len(dni) != 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DNI must be 8 digits"})
		return
	}

	client, err := h.service.GetClientByDNI(dni)
	if err != nil {
		if err.Error() == "client not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *IAMHandler) UpdateClient(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.service.UpdateClient(id, req)
	if err != nil {
		if err.Error() == "client not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// GetAllClients godoc
// @Summary Get all clients
// @Description Get a list of all clients
// @Tags clients
// @Produce json
// @Success 200 {array} Client
// @Failure 500 {object} map[string]interface{}
// @Router /clients [get]
func (h *IAMHandler) GetAllClients(c *gin.Context) {
	clients, err := h.service.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clients)
}

// ValidateRENIECByDNI godoc
// @Summary Validate DNI with RENIEC (for chatbot)
// @Description Validate if a DNI exists in RENIEC before registration
// @Tags reniec
// @Accept json
// @Produce json
// @Param dni path string true "DNI to validate"
// @Success 200 {object} ReniecValidationResult
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reniec/validate/{dni} [get]
func (h *IAMHandler) ValidateRENIECByDNI(c *gin.Context) {
	dni := c.Param("dni")
	if len(dni) != 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DNI must be 8 digits"})
		return
	}

	result, err := h.service.ValidateWithRENIEC(dni)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}