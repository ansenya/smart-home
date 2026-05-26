package repositories

import (
	"devices-api/internal/models"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DevicesRepository interface {
	List(userID uuid.UUID) ([]models.DeviceWithRelations, error)
	Get(userID, deviceID uuid.UUID) (*models.DeviceWithRelations, error)
	Update(userID, deviceID uuid.UUID, req *models.UpdateDeviceRequest) (*models.DeviceWithRelations, error)
	Delete(userID, deviceID uuid.UUID) error
	UpdateCapabilityState(userID, deviceID uuid.UUID, capType, instance string, value json.RawMessage) error
}

type devicesRepository struct {
	db *gorm.DB
}

func NewDevicesRepository(db *gorm.DB) DevicesRepository {
	return &devicesRepository{db: db}
}

func (r *devicesRepository) List(userID uuid.UUID) ([]models.DeviceWithRelations, error) {
	var devices []models.DeviceWithRelations
	err := r.db.
		Preload("Capabilities").
		Preload("Properties").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&devices).Error
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	return devices, nil
}

func (r *devicesRepository) Get(userID, deviceID uuid.UUID) (*models.DeviceWithRelations, error) {
	var device models.DeviceWithRelations
	err := r.db.
		Preload("Capabilities").
		Preload("Properties").
		Where("id = ? AND user_id = ?", deviceID, userID).
		First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *devicesRepository) Update(userID, deviceID uuid.UUID, req *models.UpdateDeviceRequest) (*models.DeviceWithRelations, error) {
	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Room != nil {
		updates["room"] = *req.Room
	}
	if len(updates) == 0 {
		return r.Get(userID, deviceID)
	}

	res := r.db.Model(&models.DeviceWithRelations{}).
		Where("id = ? AND user_id = ?", deviceID, userID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.Get(userID, deviceID)
}

func (r *devicesRepository) Delete(userID, deviceID uuid.UUID) error {
	res := r.db.Where("id = ? AND user_id = ?", deviceID, userID).
		Delete(&models.DeviceWithRelations{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *devicesRepository) UpdateCapabilityState(userID, deviceID uuid.UUID, capType, instance string, value json.RawMessage) error {
	// First verify device ownership
	var device models.DeviceWithRelations
	if err := r.db.Select("id").Where("id = ? AND user_id = ?", deviceID, userID).
		First(&device).Error; err != nil {
		return err
	}

	var cap models.Capability
	if err := r.db.Where("device_id = ? AND type = ?", deviceID, capType).
		First(&cap).Error; err != nil {
		return fmt.Errorf("capability %s not found on device %s", capType, deviceID)
	}

	existing := map[string]json.RawMessage{}
	if len(cap.State) > 0 {
		_ = json.Unmarshal(cap.State, &existing)
	}
	if instance != "" {
		instanceRaw, err := json.Marshal(instance)
		if err == nil {
			existing["instance"] = instanceRaw
		}
	}
	existing["value"] = value

	newState, err := json.Marshal(existing)
	if err != nil {
		return fmt.Errorf("marshal capability state: %w", err)
	}

	res := r.db.Model(&models.Capability{}).
		Where("id = ?", cap.ID).
		Update("state", newState)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("capability state not updated")
	}
	return nil
}
