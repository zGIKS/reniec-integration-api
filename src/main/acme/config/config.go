package config

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	RENIEC   ReniecConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Host string
	Mode string // gin mode: debug, release, test
}

type ReniecConfig struct {
	APIKey  string
	BaseURL string
}

type AppConfig struct {
	Environment string // development, production, testing
	LogLevel    string
	EnableCORS  bool
}

func (db DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.Username, db.Password, db.Name, db.SSLMode)
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "acme"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			EnableCORS:  getBoolEnv("ENABLE_CORS", true),
		},
		RENIEC: ReniecConfig{
			APIKey:  getEnv("RENIEC_API_KEY", ""),
			BaseURL: getEnv("RENIEC_BASE_URL", ""),
		},
	}

	file, err := os.Open("resources/app.properties")
	if err != nil {
		return nil, fmt.Errorf("error opening app.properties: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := expandEnvVars(strings.TrimSpace(parts[1]))

		switch key {
		case "spring.datasource.url":
			// Parse PostgreSQL URL: jdbc:postgresql://host:port/database
			if strings.Contains(value, "postgresql://") {
				// Extract from jdbc:postgresql://host:port/database
				url := strings.TrimPrefix(value, "jdbc:postgresql://")
				parts := strings.Split(url, "/")
				if len(parts) >= 2 {
					config.Database.Name = parts[1]
					hostPort := parts[0]
					if strings.Contains(hostPort, ":") {
						hostPortParts := strings.Split(hostPort, ":")
						config.Database.Host = hostPortParts[0]
						config.Database.Port = hostPortParts[1]
					} else {
						config.Database.Host = hostPort
						config.Database.Port = "5432" // default PostgreSQL port
					}
				}
			}
		case "spring.datasource.username":
			config.Database.Username = value
		case "spring.datasource.password":
			config.Database.Password = value
		case "db.sslmode":
			config.Database.SSLMode = value
		case "server.port":
			config.Server.Port = value
		case "server.host":
			config.Server.Host = value
		case "gin.mode":
			config.Server.Mode = value
		case "app.environment":
			config.App.Environment = value
		case "app.log.level":
			config.App.LogLevel = value
		case "app.cors.enabled":
			if enabled, err := strconv.ParseBool(value); err == nil {
				config.App.EnableCORS = enabled
			}
		case "reniec.ruc.api.key":
			config.RENIEC.APIKey = value
		case "reniec.ruc.api.base.url":
			config.RENIEC.BaseURL = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading app.properties: %w", err)
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// expandEnvVars expands environment variables in the format ${VAR_NAME}
func expandEnvVars(value string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(value, func(match string) string {
		// Extract variable name (remove ${ and })
		varName := match[2 : len(match)-1]
		return os.Getenv(varName)
	})
}
