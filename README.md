# ACME Backend API

> **Enterprise Management System** - RESTful API for appointment scheduling, service catalog, and client management with RENIEC integration

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Docker Setup](#docker-setup)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Deployment](#deployment)

## Quick Start

### Prerequisites

- **Go 1.21+**
- **PostgreSQL 15+**
- **Docker & Docker Compose** (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd backend
   ```

2. **Install dependencies**
   ```bash
   cd src/main/acme
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

5. **Access the API**
   - **API Base URL:** `http://localhost:8080/api/v1`
   - **Swagger Documentation:** `http://localhost:8080/swagger/index.html`
   - **Health Check:** `http://localhost:8080/health`

## Architecture

The ACME Backend follows a **clean architecture** pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │───▶│  Business Logic │───▶│  Data Access    │
│   (Handlers)    │    │   (Services)    │    │  (Repository)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                        │                        │
         ▼                        ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Gin Router   │    │   Domain Models │    │   PostgreSQL    │
│   & Middleware  │    │   & Validation  │    │    Database     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Core Modules

| Module | Description | Features |
|--------|-------------|----------|
| **Appointments** | Appointment scheduling system | Create, cancel, availability check |
| **IAM (Identity)** | Client management with RENIEC | Client CRUD, DNI validation |
| **Catalog** | Service catalog management | Service CRUD, pricing |
| **Employees** | Employee management | Employee directory |
| **Audit** | System audit logging | Activity tracking |

## API Documentation

### Interactive Documentation
- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **API Version:** `v1.0`
- **Base Path:** `/api/v1`

### Authentication & Security
- **CORS:** Configurable cross-origin support
- **Environment-based:** Development/Production modes
- **Input Validation:** Comprehensive request validation

## Configuration

### Environment Variables

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=acme
DB_USER=postgres
DB_PASSWORD=your_password
DB_SSLMODE=disable

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=debug

# Application Configuration
APP_ENV=development
LOG_LEVEL=info
ENABLE_CORS=true

# RENIEC API Integration
RENIEC_API_KEY=your_api_key
RENIEC_BASE_URL=https://api.reniec.gob.pe/v1
```

### Application Properties

The system also supports Java-style properties files for additional configuration in `src/main/resources/app.properties`.

## Docker Setup

### Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Manual Docker Build

```bash
# Build the image
docker build -t acme-backend .

# Run with local PostgreSQL
docker run --network host --env-file .env acme-backend

# Run with port mapping
docker run --env-file .env -p 8080:8080 acme-backend
```

## Project Structure

```
backend/
├── README.md
├── Dockerfile
├── docker-compose.yml
├── .env
├── .gitignore
└── src/main/
    ├── acme/
    │   ├── appointments/        # Appointment management
    │   ├── audit/              # Audit logging
    │   ├── catalog/            # Service catalog
    │   ├── config/             # Configuration management
    │   ├── database/           # Database connection & migrations
    │   ├── docs/               # Swagger documentation
    │   ├── employees/          # Employee management
    │   ├── iam/                # Identity & Access Management
    │   ├── router/             # HTTP routing
    │   ├── go.mod              # Go dependencies
    │   └── main.go             # Application entry point
    └── resources/
        └── app.properties      # Application configuration
```

## API Endpoints

### Appointments Management

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `POST` | `/appointments` | Create new appointment | `CreateAppointmentRequest` |
| `GET` | `/appointments/{id}` | Get appointment by ID | - |
| `GET` | `/appointments/{id}/details` | Get appointment with full details | - |
| `PUT` | `/appointments/{id}` | Update appointment | `UpdateAppointmentRequest` |
| `PUT` | `/appointments/{id}/cancel` | Cancel appointment (admin) | `CancelAppointmentRequest` |
| `PUT` | `/appointments/{id}/cancel-by-client` | Cancel appointment (client) | `{dni, reason}` |
| `PUT` | `/appointments/{id}/cancel-by-employee` | Cancel appointment (employee) | `{email, reason}` |
| `GET` | `/appointments/date-range` | Get appointments by date range | `?start_date&end_date` |
| `GET` | `/appointments/client/{client_id}` | Get client's appointments | - |
| `GET` | `/appointments/availability` | Check availability | `?date&time&service_id` |

### Client Management (IAM)

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `POST` | `/clients` | Create new client | `CreateClientRequest` |
| `GET` | `/clients` | Get all clients | - |
| `GET` | `/clients/{id}` | Get client by ID | - |
| `GET` | `/clients/dni/{dni}` | Get client by DNI | - |
| `PUT` | `/clients/{id}` | Update client | `UpdateClientRequest` |

### Service Catalog

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `POST` | `/services` | Create new service | `CreateServiceRequest` |
| `GET` | `/services` | Get all services | - |
| `GET` | `/services/{id}` | Get service by ID | - |
| `PUT` | `/services/{id}` | Update service | `UpdateServiceRequest` |
| `DELETE` | `/services/{id}` | Delete service | - |
| `GET` | `/services/price-range` | Get services by price range | `?min&max` |

### Employee Management

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `GET` | `/employees` | Get all employees | - |
| `GET` | `/employees/{id}` | Get employee by ID | - |

### RENIEC Integration

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `GET` | `/reniec/validate/{dni}` | Validate DNI with RENIEC | - |

### System Health

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/health` | Health check | `{"status": "healthy"}` |

## Database Schema

### Core Entities

```sql
-- Clients
CREATE TABLE clients (
    id UUID PRIMARY KEY,
    dni VARCHAR(8) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    second_last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    reniec_validated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Services
CREATE TABLE services (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    benefits TEXT,
    includes TEXT,
    contraindications TEXT,
    recommended_frequency VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Appointments
CREATE TABLE appointments (
    id UUID PRIMARY KEY,
    client_id UUID REFERENCES clients(id),
    service_id UUID REFERENCES services(id),
    appointment_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled',
    attended_by VARCHAR(255),
    cancellation_reason TEXT,
    cancelled_by VARCHAR(255),
    cancelled_by_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Employees
CREATE TABLE employees (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    paternal_surname VARCHAR(100) NOT NULL,
    maternal_surname VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### API Testing

Use the provided Swagger UI at `http://localhost:8080/swagger/index.html` to test all endpoints interactively.

### Example Requests

#### Create a Client
```bash
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Content-Type: application/json" \
  -d '{
    "dni": "12345678",
    "first_name": "Juan",
    "last_name": "Perez",
    "email": "juan.perez@email.com",
    "phone": "+51987654321"
  }'
```

#### Create an Appointment
```bash
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "uuid-here",
    "service_id": "uuid-here",
    "appointment_date": "2024-01-15",
    "start_time": "10:00"
  }'
```

## Deployment

### Production Checklist

- [ ] Set `GIN_MODE=release`
- [ ] Configure production database
- [ ] Set up SSL/TLS certificates
- [ ] Configure reverse proxy (Nginx/Apache)
- [ ] Set up monitoring and logging
- [ ] Configure backup strategies
- [ ] Update RENIEC API credentials

### Environment Setup

```bash
# Production environment variables
export APP_ENV=production
export GIN_MODE=release
export LOG_LEVEL=error
export DB_SSLMODE=require
```

### Health Monitoring

The API provides a health check endpoint at `/health` for load balancers and monitoring systems.

---

## Support & Contact

- **API Documentation:** [Swagger UI](http://localhost:8080/swagger/index.html)
- **Health Status:** [Health Check](http://localhost:8080/health)

---

**Built with Go, Gin, PostgreSQL, and Docker**
