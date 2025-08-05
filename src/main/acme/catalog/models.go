package catalog

import (
	"time"
)

type Service struct {
	ID                   string    `json:"id" db:"id"`
	Name                 string    `json:"name" db:"name"`
	Price                float64   `json:"price" db:"price"`
	DurationMinutes      int       `json:"duration_minutes" db:"duration_minutes"`
	Description          *string   `json:"description" db:"description"`
	Benefits             *string   `json:"benefits" db:"benefits"`
	RecommendedFrequency *string   `json:"recommended_frequency" db:"recommended_frequency"`
	Includes             *string   `json:"includes" db:"includes"`
	Contraindications    *string   `json:"contraindications" db:"contraindications"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type CreateServiceRequest struct {
	Name                 string  `json:"name" binding:"required"`
	Price                float64 `json:"price" binding:"required,min=0"`
	DurationMinutes      int     `json:"duration_minutes" binding:"required,min=1"`
	Description          *string `json:"description"`
	Benefits             *string `json:"benefits"`
	RecommendedFrequency *string `json:"recommended_frequency"`
	Includes             *string `json:"includes"`
	Contraindications    *string `json:"contraindications"`
}

type UpdateServiceRequest struct {
	Name                 *string  `json:"name"`
	Price                *float64 `json:"price"`
	DurationMinutes      *int     `json:"duration_minutes"`
	Description          *string  `json:"description"`
	Benefits             *string  `json:"benefits"`
	RecommendedFrequency *string  `json:"recommended_frequency"`
	Includes             *string  `json:"includes"`
	Contraindications    *string  `json:"contraindications"`
}