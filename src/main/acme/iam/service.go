package iam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"acme/config"
)

type IAMService struct {
	repo   *Repository
	config *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *IAMService {
	return &IAMService{
		repo:   repo,
		config: cfg,
	}
}

func (s *IAMService) ValidateWithRENIEC(dni string) (*ReniecValidationResult, error) {
	url := fmt.Sprintf("%s/v1/reniec/dni?numero=%s", s.config.RENIEC.BaseURL, dni)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &ReniecValidationResult{
			IsValid: false,
			Error:   "Error creating request",
		}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.RENIEC.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &ReniecValidationResult{
			IsValid: false,
			Error:   "Error making request to RENIEC",
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &ReniecValidationResult{
			IsValid: false,
			Error:   fmt.Sprintf("RENIEC API returned status: %d", resp.StatusCode),
		}, nil
	}

	var reniecData ReniecResponse
	if err := json.NewDecoder(resp.Body).Decode(&reniecData); err != nil {
		return &ReniecValidationResult{
			IsValid: false,
			Error:   "Error decoding RENIEC response",
		}, err
	}

	return &ReniecValidationResult{
		IsValid: true,
		Data:    reniecData,
	}, nil
}

func (s *IAMService) CreateClient(req CreateClientRequest) (*Client, error) {
	// Primero validar con RENIEC que el DNI existe
	reniecResult, err := s.ValidateWithRENIEC(req.DNI)
	if err != nil {
		return nil, fmt.Errorf("error validating with RENIEC: %w", err)
	}

	// Solo permitir registro si existe en RENIEC
	if !reniecResult.IsValid {
		return nil, fmt.Errorf("DNI not found in RENIEC. Client registration not allowed")
	}

	// Verificar si el cliente ya existe
	existingClient, _ := s.repo.GetClientByDNI(req.DNI)
	if existingClient != nil {
		return nil, fmt.Errorf("client with DNI %s already exists", req.DNI)
	}

	client := &Client{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		SecondLastName:  req.SecondLastName,
		DNI:             req.DNI,
		Email:           req.Email,
		Phone:           req.Phone,
		ReniecValidated: true,
	}

	// Validar que los datos coincidan con RENIEC
	if !s.validateClientDataWithRENIEC(req, reniecResult.Data) {
		return nil, fmt.Errorf("client data does not match RENIEC records")
	}

	if err := s.repo.CreateClient(client); err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return client, nil
}

func (s *IAMService) validateClientDataWithRENIEC(clientReq CreateClientRequest, reniecData ReniecResponse) bool {
	clientFullName := strings.ToUpper(strings.TrimSpace(fmt.Sprintf("%s %s", clientReq.FirstName, clientReq.LastName)))
	if clientReq.SecondLastName != nil && *clientReq.SecondLastName != "" {
		clientFullName += " " + strings.ToUpper(strings.TrimSpace(*clientReq.SecondLastName))
	}

	reniecFullName := strings.ToUpper(strings.TrimSpace(reniecData.FullName))

	return clientFullName == reniecFullName && clientReq.DNI == reniecData.DocumentNumber
}

func (s *IAMService) GetClientByID(id string) (*Client, error) {
	return s.repo.GetClientByID(id)
}

func (s *IAMService) GetClientByDNI(dni string) (*Client, error) {
	return s.repo.GetClientByDNI(dni)
}

func (s *IAMService) UpdateClient(id string, req UpdateClientRequest) (*Client, error) {
	if err := s.repo.UpdateClient(id, req); err != nil {
		return nil, fmt.Errorf("error updating client: %w", err)
	}

	return s.repo.GetClientByID(id)
}

func (s *IAMService) GetAllClients() ([]Client, error) {
	return s.repo.GetAllClients()
}