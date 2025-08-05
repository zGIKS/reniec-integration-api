package catalog

import (
	"database/sql"

	"acme/appointments"
	"acme/audit"
	"acme/config"
	"acme/employees"
	"acme/iam"
)

// ServiceFactory implements the Factory pattern for creating services
type ServiceFactory struct {
	db     *sql.DB
	config *config.Config
}

func NewServiceFactory(db *sql.DB, cfg *config.Config) *ServiceFactory {
	return &ServiceFactory{db: db, config: cfg}
}

// CreateServices creates all application services using dependency injection
func (f *ServiceFactory) CreateServices() (*AppServices, error) {
	// Create repositories
	auditRepo := audit.NewRepository(f.db)
	iamRepo := iam.NewRepository(f.db)
	catalogRepo := NewRepository(f.db)
	appointmentsRepo := appointments.NewRepository(f.db)
	employeesRepo := employees.NewRepository(f.db)

	// Create services with dependencies
	auditService := audit.NewService(auditRepo)
	iamService := iam.NewService(iamRepo, f.config)
	catalogService := NewService(catalogRepo)
	appointmentsService := appointments.NewService(appointmentsRepo, auditService)
	employeesService := employees.NewService(employeesRepo)

	return &AppServices{
		Audit:        auditService,
		IAM:          iamService,
		Catalog:      catalogService,
		Appointments: appointmentsService,
		Employees:    employeesService,
	}, nil
}

// CreateHandlers creates all HTTP handlers
func (f *ServiceFactory) CreateHandlers(services *AppServices) *AppHandlers {
	return &AppHandlers{
		IAM:          iam.NewIAMHandler(services.IAM),
		Catalog:      NewCatalogHandler(services.Catalog),
		Appointments: appointments.NewAppointmentsHandler(services.Appointments),
		Employees:    employees.NewEmployeesHandler(services.Employees),
	}
}

// AppServices holds all application services
type AppServices struct {
	Audit        *audit.Service
	IAM          *iam.IAMService
	Catalog      *CatalogService
	Appointments *appointments.AppointmentService
	Employees    *employees.EmployeeService
}

// AppHandlers holds all HTTP handlers
type AppHandlers struct {
	IAM          *iam.IAMHandler
	Catalog      *CatalogHandler
	Appointments *appointments.AppointmentsHandler
	Employees    *employees.EmployeesHandler
}
