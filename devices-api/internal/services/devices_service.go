package services

import (
	"devices-api/internal/models"
	"devices-api/internal/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
)

type DevicesService interface {
	List(userID uuid.UUID) ([]models.DeviceWithRelations, error)
	Get(userID, deviceID uuid.UUID) (*models.DeviceWithRelations, error)
	Update(userID, deviceID uuid.UUID, req *models.UpdateDeviceRequest) (*models.DeviceWithRelations, error)
	Delete(userID, deviceID uuid.UUID) error
	SetCapabilityValue(userID, deviceID uuid.UUID, capType string, req *models.SetCapabilityRequest) error
}

type devicesService struct {
	repo   repositories.DevicesRepository
	mqtt   mqtt.Client
	log    *slog.Logger
}

func NewDevicesService(repo repositories.DevicesRepository, mqttClient mqtt.Client, log *slog.Logger) DevicesService {
	return &devicesService{repo: repo, mqtt: mqttClient, log: log}
}

func (s *devicesService) List(userID uuid.UUID) ([]models.DeviceWithRelations, error) {
	return s.repo.List(userID)
}

func (s *devicesService) Get(userID, deviceID uuid.UUID) (*models.DeviceWithRelations, error) {
	d, err := s.repo.Get(userID, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return d, nil
}

func (s *devicesService) Update(userID, deviceID uuid.UUID, req *models.UpdateDeviceRequest) (*models.DeviceWithRelations, error) {
	d, err := s.repo.Update(userID, deviceID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return d, nil
}

func (s *devicesService) Delete(userID, deviceID uuid.UUID) error {
	if err := s.repo.Delete(userID, deviceID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDeviceNotFound
		}
		return err
	}
	return nil
}

func (s *devicesService) SetCapabilityValue(userID, deviceID uuid.UUID, capType string, req *models.SetCapabilityRequest) error {
	if s.mqtt == nil || !s.mqtt.IsConnected() {
		return fmt.Errorf("mqtt not connected")
	}

	if err := s.repo.UpdateCapabilityState(userID, deviceID, capType, req.Instance, req.Value); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDeviceNotFound
		}
		return err
	}

	topic := fmt.Sprintf("%s/%s/capabilities/%s/set", userID.String(), deviceID.String(), capType)
	payload, err := json.Marshal(map[string]any{
		"instance": req.Instance,
		"value":    req.Value,
		"ts":       time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	token := s.mqtt.Publish(topic, 0, false, payload)
	if !token.WaitTimeout(2 * time.Second) {
		return fmt.Errorf("mqtt publish timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt publish: %w", err)
	}

	s.log.Info("capability published",
		slog.String("topic", topic),
		slog.String("payload", string(payload)),
	)
	return nil
}
