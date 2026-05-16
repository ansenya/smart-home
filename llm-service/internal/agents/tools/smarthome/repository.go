package smarthome

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type deviceRepo struct {
	db *gorm.DB
}

func newDeviceRepo(db *gorm.DB) *deviceRepo {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) listDevices(ctx context.Context, userID string) ([]Device, error) {
	var devices []Device
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	return devices, nil
}

func (r *deviceRepo) getDevice(ctx context.Context, deviceID, userID string) (*Device, error) {
	var device Device
	err := r.db.WithContext(ctx).
		Preload("Capabilities").
		Preload("Properties").
		Where("id = ? AND user_id = ?", deviceID, userID).
		First(&device).Error
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	return &device, nil
}

func (r *deviceRepo) updateCapabilityState(ctx context.Context, deviceID, capType string, value json.RawMessage) error {
	// Read current state to preserve the instance field
	var cap Capability
	if err := r.db.WithContext(ctx).
		Where("device_id = ? AND type = ?", deviceID, capType).
		First(&cap).Error; err != nil {
		return fmt.Errorf("capability %s not found on device %s", capType, deviceID)
	}

	// Parse existing state to keep instance
	var existing map[string]json.RawMessage
	if len(cap.State) > 0 {
		_ = json.Unmarshal(cap.State, &existing)
	}
	if existing == nil {
		existing = make(map[string]json.RawMessage)
	}
	existing["value"] = value

	newState, err := json.Marshal(existing)
	if err != nil {
		return err
	}

	res := r.db.WithContext(ctx).Model(&Capability{}).
		Where("device_id = ? AND type = ?", deviceID, capType).
		Update("state", string(newState))
	if res.Error != nil {
		return fmt.Errorf("update capability state: %w", res.Error)
	}
	return nil
}
