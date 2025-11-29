package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/repositories"
	"github.com/streadway/amqp"
)

type TenantService struct {
	repo    *repositories.TenantRepository
	channel *amqp.Channel
}

func NewTenantService(repo *repositories.TenantRepository, channel *amqp.Channel) *TenantService {
	return &TenantService{
		repo:    repo,
		channel: channel,
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, req dto.CreateTenantRequest) (*dto.TenantResponse, error) {

	//Convertir a minusculas
	req.TenantID = strings.ToLower(req.TenantID)
	//Validar que sea solo una palabra el ID
	if !haveOneWord(req.TenantID) {
		return nil, fmt.Errorf("TenantID debe ser una sola palabra")
	}

	// Guardar en BD
	tenant, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Publicar evento
	s.publishEvent("tenant.created", map[string]interface{}{
		"tenant_id":   req.TenantID,
		"tenant_name": req.TenantName,
	})

	return tenant, nil
}

func (s *TenantService) GetTenant(tenantID string) (*dto.TenantResponse, error) {
	return s.repo.FindByID(tenantID)
}

func (s *TenantService) GetAllTenants() ([]dto.TenantResponse, error) {
	return s.repo.GetAll()
}

func (s *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	// Eliminar de BD
	if err := s.repo.Delete(tenantID); err != nil {
		return err
	}

	// Publicar evento
	s.publishEvent("tenant.deleted", map[string]interface{}{
		"tenant_id": tenantID,
	})

	return nil
}

func (s *TenantService) publishEvent(eventType string, data map[string]interface{}) {
	if s.channel == nil {
		return
	}

	eventPayload := map[string]interface{}{
		"event_type": eventType,
		"data":       data,
	}

	body, err := json.Marshal(eventPayload)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	s.channel.Publish(
		"tenants_events",
		eventType,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	fmt.Println("Tenant published")
}

func haveOneWord(cadena string) bool {
	return len(strings.TrimSpace(cadena)) > 0 && !strings.Contains(cadena, " ")
}
