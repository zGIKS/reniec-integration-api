package employees

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployeesHandler struct {
	service *EmployeeService
}

func NewEmployeesHandler(service *EmployeeService) *EmployeesHandler {
	return &EmployeesHandler{service: service}
}

// GetAllEmployees godoc
// @Summary Get all employees
// @Description Get a list of all employees
// @Tags employees
// @Produce json
// @Success 200 {array} Employee
// @Failure 500 {object} map[string]interface{}
// @Router /employees [get]
func (h *EmployeesHandler) GetAllEmployees(c *gin.Context) {
	employees, err := h.service.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}


// GetEmployeeByID godoc
// @Summary Get employee by ID
// @Description Get employee details by ID
// @Tags employees
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} Employee
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /employees/{id} [get]
func (h *EmployeesHandler) GetEmployeeByID(c *gin.Context) {
	id := c.Param("id")

	employee, err := h.service.GetEmployeeByID(id)
	if err != nil {
		if err.Error() == "employee not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}