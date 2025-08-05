package database

import (
	"database/sql"
	"fmt"
)

func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		
		`CREATE TABLE IF NOT EXISTS clients (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			second_last_name VARCHAR(100),
			dni VARCHAR(8) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20),
			registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			reniec_validated BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS services (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			duration_minutes INTEGER NOT NULL,
			description TEXT,
			benefits TEXT,
			recommended_frequency VARCHAR(255),
			includes TEXT,
			contraindications TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS employees (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(100) NOT NULL,
			paternal_surname VARCHAR(100) NOT NULL,
			maternal_surname VARCHAR(100),
			role VARCHAR(50) DEFAULT 'specialist',
			phone VARCHAR(20),
			email VARCHAR(255) UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS audit_logs (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			table_name VARCHAR(50) NOT NULL,
			record_id UUID NOT NULL,
			action VARCHAR(20) NOT NULL CHECK (action IN ('CREATE', 'UPDATE', 'DELETE', 'CANCEL')),
			old_values JSONB,
			new_values JSONB,
			changed_by VARCHAR(100),
			changed_by_type VARCHAR(20) CHECK (changed_by_type IN ('client', 'employee', 'system')),
			reason TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS appointments (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
			service_id UUID NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
			appointment_date DATE NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			attended_by UUID REFERENCES employees(id),
			status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
			cancelled_by VARCHAR(100),
			cancelled_by_type VARCHAR(20) CHECK (cancelled_by_type IN ('client', 'employee')),
			cancellation_reason TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(appointment_date, start_time, attended_by)
		)`,

		`CREATE INDEX IF NOT EXISTS idx_clients_dni ON clients(dni)`,
		`CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email)`,
		`CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_table_record ON audit_logs(table_name, record_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_appointments_date ON appointments(appointment_date)`,
		`CREATE INDEX IF NOT EXISTS idx_appointments_client ON appointments(client_id)`,
		`CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status)`,
		`CREATE INDEX IF NOT EXISTS idx_appointments_attended_by ON appointments(attended_by)`,

		`CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql'`,

		`DROP TRIGGER IF EXISTS update_clients_updated_at ON clients`,
		`CREATE TRIGGER update_clients_updated_at BEFORE UPDATE ON clients FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,

		`DROP TRIGGER IF EXISTS update_services_updated_at ON services`,
		`CREATE TRIGGER update_services_updated_at BEFORE UPDATE ON services FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,

		`DROP TRIGGER IF EXISTS update_employees_updated_at ON employees`,
		`CREATE TRIGGER update_employees_updated_at BEFORE UPDATE ON employees FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,

		`DROP TRIGGER IF EXISTS update_appointments_updated_at ON appointments`,
		`CREATE TRIGGER update_appointments_updated_at BEFORE UPDATE ON appointments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error executing query: %s\nError: %w", query, err)
		}
	}

	return nil
}

func SeedData(db *sql.DB) error {
	checkQuery := `SELECT COUNT(*) FROM services`
	var count int
	err := db.QueryRow(checkQuery).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking if services exist: %w", err)
	}

	if count > 0 {
		return nil
	}

	// Seed employees first
	employeeQueries := []string{
		`INSERT INTO employees (name, paternal_surname, maternal_surname, phone, email) VALUES 
		('María', 'Fernández', 'Silva', '999111222', 'maria.fernandez@acme.com'),
		('José', 'Mendoza', 'Torres', '999333444', 'jose.mendoza@acme.com'),
		('Carmen', 'Vargas', 'Ruiz', '999555666', 'carmen.vargas@acme.com'),
		('Ana', 'Jiménez', 'Castro', '999777888', 'ana.jimenez@acme.com'),
		('Luis', 'Morales', 'Vega', '999999000', 'luis.morales@acme.com')`,
	}

	for _, query := range employeeQueries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error seeding employees: %w", err)
		}
	}

	seedQueries := []string{
		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Limpieza Facial Profunda', 80.00, 60, 'Limpieza facial completa con extracción de puntos negros y mascarilla hidratante', 'Elimina impurezas, destapa poros, mejora textura de la piel', 'Cada 4-6 semanas', 'Desmaquillado, limpieza, exfoliación, extracción, mascarilla, hidratación', 'Acné activo severo, heridas abiertas, rosácea severa')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Facial Anti-edad', 120.00, 75, 'Tratamiento facial con productos anti-envejecimiento y técnicas de lifting', 'Reduce líneas de expresión, mejora firmeza, hidrata profundamente', 'Cada 3-4 semanas', 'Limpieza, sérum anti-edad, masaje facial, mascarilla colágeno', 'Embarazo, alergia a ingredientes activos')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Masaje Relajante Corporal', 100.00, 60, 'Masaje corporal con aceites esenciales para relajación total', 'Reduce estrés, mejora circulación, alivia tensión muscular', 'Según necesidad', 'Masaje con aceites aromáticos, música relajante', 'Lesiones recientes, fiebre, embarazo primer trimestre')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Tratamiento Corporal Reductivo', 150.00, 90, 'Tratamiento para reducir medidas con técnicas de drenaje y productos específicos', 'Reduce medidas, mejora circulación, tonifica la piel', 'Semanal por 8 sesiones', 'Exfoliación, masaje reductivo, cremas específicas, vendas', 'Embarazo, problemas cardiovasculares, varices severas')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Manicure Completo', 35.00, 45, 'Cuidado completo de manos y uñas con esmaltado', 'Manos suaves, uñas bien cuidadas y esmaltadas', 'Cada 2-3 semanas', 'Limado, cutícula, exfoliación, hidratación, esmaltado', 'Infecciones en uñas o manos')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Pedicure Completo', 45.00, 60, 'Cuidado completo de pies y uñas con esmaltado', 'Pies suaves, uñas cuidadas, relajación', 'Cada 3-4 semanas', 'Remojo, limado, cutícula, exfoliación, masaje, esmaltado', 'Heridas abiertas, hongos activos')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Facial Express', 50.00, 30, 'Limpieza facial rápida ideal para mantenimiento', 'Limpieza básica, hidratación inmediata', 'Cada 2 semanas', 'Limpieza, tónico, mascarilla express, hidratante', 'Acné severo')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Depilación con Cera Piernas Completas', 60.00, 45, 'Depilación con cera caliente para piernas completas', 'Piel suave por 3-4 semanas, cabello más fino', 'Cada 4-6 semanas', 'Limpieza, aplicación de cera, crema calmante', 'Piel irritada, quemaduras solares recientes')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Microdermoabrasión', 95.00, 50, 'Exfoliación profunda con cristales para renovación celular', 'Mejora textura, reduce manchas, estimula renovación', 'Cada 2-3 semanas', 'Limpieza, microdermoabrasión, mascarilla calmante', 'Piel muy sensible, rosácea, heridas')`,

		`INSERT INTO services (name, price, duration_minutes, description, benefits, recommended_frequency, includes, contraindications) VALUES 
		('Radiofrecuencia Facial', 180.00, 60, 'Tratamiento con radiofrecuencia para firmeza y lifting', 'Efecto lifting, mejora firmeza, estimula colágeno', 'Cada 2 semanas por 6 sesiones', 'Limpieza, gel conductor, aplicación RF, hidratación', 'Embarazo, marcapasos, implantes metálicos')`,
	}

	for _, query := range seedQueries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error seeding data: %w", err)
		}
	}

	return nil
}