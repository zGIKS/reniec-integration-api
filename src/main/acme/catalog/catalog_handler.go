package catalog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	service *CatalogService
}

func NewCatalogHandler(service *CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

func (h *CatalogHandler) CreateService(c *gin.Context) {
	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.service.CreateService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (h *CatalogHandler) GetServiceByID(c *gin.Context) {
	id := c.Param("id")

	service, err := h.service.GetServiceByID(id)
	if err != nil {
		if err.Error() == "service not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// GetAllServices godoc
// @Summary Get all services
// @Description Get a list of all available services
// @Tags services
// @Produce json
// @Success 200 {array} Service
// @Failure 500 {object} map[string]interface{}
// @Router /services [get]
func (h *CatalogHandler) GetAllServices(c *gin.Context) {
	services, err := h.service.GetAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}


func (h *CatalogHandler) UpdateService(c *gin.Context) {
	id := c.Param("id")

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.service.UpdateService(id, req)
	if err != nil {
		if err.Error() == "service not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

func (h *CatalogHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteService(id)
	if err != nil {
		if err.Error() == "service not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *CatalogHandler) GetServicesByPriceRange(c *gin.Context) {
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	if minPriceStr == "" || maxPriceStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price and max_price query parameters are required"})
		return
	}

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price value"})
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price value"})
		return
	}

	services, err := h.service.GetServicesByPriceRange(minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

