package router

import (
	"acme/config"
	"acme/catalog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, handlers *catalog.AppHandlers, cfg *config.Config) {

	// Add CORS middleware if enabled
	if cfg.App.EnableCORS {
		r.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			
			c.Next()
		})
	}

	api := r.Group("/api/v1")
	{
		// RENIEC validation endpoint (for chatbot flow)
		reniec := api.Group("/reniec")
		{
			reniec.GET("/validate/:dni", handlers.IAM.ValidateRENIECByDNI)
		}

		clients := api.Group("/clients")
		{
			clients.POST("", handlers.IAM.CreateClient)
			clients.GET("", handlers.IAM.GetAllClients)
			clients.GET("/:id", handlers.IAM.GetClientByID)
			clients.PUT("/:id", handlers.IAM.UpdateClient)
			clients.GET("/dni/:dni", handlers.IAM.GetClientByDNI)
		}

		services := api.Group("/services")
		{
			services.POST("", handlers.Catalog.CreateService)
			services.GET("", handlers.Catalog.GetAllServices)
			services.GET("/:id", handlers.Catalog.GetServiceByID)
			services.PUT("/:id", handlers.Catalog.UpdateService)
			services.DELETE("/:id", handlers.Catalog.DeleteService)
			services.GET("/price-range", handlers.Catalog.GetServicesByPriceRange)
		}

		employees := api.Group("/employees")
		{
			employees.GET("", handlers.Employees.GetAllEmployees)
			employees.GET("/:id", handlers.Employees.GetEmployeeByID)
		}

		appointmentsGroup := api.Group("/appointments")
		{
			appointmentsGroup.POST("", handlers.Appointments.CreateAppointment)
			appointmentsGroup.GET("/:id", handlers.Appointments.GetAppointmentByID)
			appointmentsGroup.GET("/:id/details", handlers.Appointments.GetAppointmentWithDetails)
			appointmentsGroup.PUT("/:id", handlers.Appointments.UpdateAppointment)
			appointmentsGroup.PUT("/:id/cancel", handlers.Appointments.CancelAppointment)
			appointmentsGroup.PUT("/:id/cancel-by-client", handlers.Appointments.CancelAppointmentByClient)
			appointmentsGroup.PUT("/:id/cancel-by-employee", handlers.Appointments.CancelAppointmentByEmployee)
			appointmentsGroup.GET("/date-range", handlers.Appointments.GetAppointmentsByDateRange)
			appointmentsGroup.GET("/client/:client_id", handlers.Appointments.GetAppointmentsByClient)
			appointmentsGroup.GET("/availability", handlers.Appointments.CheckAvailability)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"message": "ACME Backend API is running",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}