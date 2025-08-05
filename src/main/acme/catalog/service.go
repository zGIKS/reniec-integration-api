package catalog

import (
	"fmt"
)

type CatalogService struct {
	repo *Repository
}

func NewService(repo *Repository) *CatalogService {
	return &CatalogService{
		repo: repo,
	}
}

func (s *CatalogService) CreateService(req CreateServiceRequest) (*Service, error) {
	service := &Service{
		Name:                 req.Name,
		Price:                req.Price,
		DurationMinutes:      req.DurationMinutes,
		Description:          req.Description,
		Benefits:             req.Benefits,
		RecommendedFrequency: req.RecommendedFrequency,
		Includes:             req.Includes,
		Contraindications:    req.Contraindications,
	}

	if err := s.repo.CreateService(service); err != nil {
		return nil, fmt.Errorf("error creating service: %w", err)
	}

	return service, nil
}

func (s *CatalogService) GetServiceByID(id string) (*Service, error) {
	return s.repo.GetServiceByID(id)
}

func (s *CatalogService) GetAllServices() ([]Service, error) {
	return s.repo.GetAllServices()
}

func (s *CatalogService) UpdateService(id string, req UpdateServiceRequest) (*Service, error) {
	if err := s.repo.UpdateService(id, req); err != nil {
		return nil, fmt.Errorf("error updating service: %w", err)
	}

	return s.repo.GetServiceByID(id)
}

func (s *CatalogService) DeleteService(id string) error {
	return s.repo.DeleteService(id)
}

func (s *CatalogService) GetServicesByPriceRange(minPrice, maxPrice float64) ([]Service, error) {
	if minPrice < 0 || maxPrice < 0 {
		return nil, fmt.Errorf("prices must be non-negative")
	}
	if minPrice > maxPrice {
		return nil, fmt.Errorf("min price cannot be greater than max price")
	}

	return s.repo.GetServicesByPriceRange(minPrice, maxPrice)
}