package smarthome

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type deviceRepo struct {
	db   *gorm.DB
	mqtt mqtt.Client
}

func newDeviceRepo(db *gorm.DB, mqttClient mqtt.Client) *deviceRepo {
	return &deviceRepo{db: db, mqtt: mqttClient}
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
	q := r.db.WithContext(ctx).
		Preload("Capabilities").
		Preload("Properties")

	// Accept either a UUID or a device name/UID — the LLM frequently passes the
	// human-readable name ("esp32-XXXX") instead of the UUID id.
	if _, err := uuid.Parse(deviceID); err == nil {
		q = q.Where("id = ? AND user_id = ?", deviceID, userID)
	} else {
		q = q.Where("name = ? AND user_id = ?", deviceID, userID)
	}

	if err := q.First(&device).Error; err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	return &device, nil
}

func (r *deviceRepo) readInstance(ctx context.Context, deviceID, capType string) (string, bool) {
	var cap Capability
	if err := r.db.WithContext(ctx).
		Where("device_id = ? AND type = ?", deviceID, capType).
		First(&cap).Error; err != nil {
		return "", false
	}
	if len(cap.State) == 0 {
		return "", false
	}
	var existing map[string]json.RawMessage
	if err := json.Unmarshal(cap.State, &existing); err != nil {
		return "", false
	}
	raw, ok := existing["instance"]
	if !ok {
		return "", false
	}
	var inst string
	if err := json.Unmarshal(raw, &inst); err != nil {
		return "", false
	}
	return inst, inst != ""
}

func (r *deviceRepo) updateCapabilityState(ctx context.Context, deviceID, capType, instance string, value json.RawMessage) error {
	var cap Capability
	if err := r.db.WithContext(ctx).
		Where("device_id = ? AND type = ?", deviceID, capType).
		First(&cap).Error; err != nil {
		return fmt.Errorf("capability %s not found on device %s", capType, deviceID)
	}

	existing := make(map[string]json.RawMessage)
	if len(cap.State) > 0 {
		_ = json.Unmarshal(cap.State, &existing)
	}
	if instance != "" {
		b, _ := json.Marshal(instance)
		existing["instance"] = b
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

func (r *deviceRepo) publishSet(userID, deviceID, capType, instance string, value json.RawMessage) error {
	if r.mqtt == nil || !r.mqtt.IsConnected() {
		return fmt.Errorf("mqtt not connected")
	}

	topic := fmt.Sprintf("%s/%s/capabilities/%s/set", userID, deviceID, capType)
	payload := map[string]any{
		"value": value,
		"ts":    time.Now().Unix(),
	}
	if instance != "" {
		payload["instance"] = instance
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	token := r.mqtt.Publish(topic, 0, false, body)
	if !token.WaitTimeout(2 * time.Second) {
		return fmt.Errorf("publish timeout")
	}
	return token.Error()
}
