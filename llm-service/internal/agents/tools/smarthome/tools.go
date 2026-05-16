package smarthome

import (
	"context"
	"encoding/json"
	"fmt"
	"llm-service/internal/agents/tools"
	"llm-service/internal/clients"
	"llm-service/internal/ctxkeys"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, reg tools.Registry) []clients.Tool {
	repo := newDeviceRepo(db)

	reg.Register("list_devices", listDevicesHandler(repo))
	reg.Register("get_device_state", getDeviceStateHandler(repo))
	reg.Register("control_device", controlDeviceHandler(repo))

	return Definitions()
}

func Definitions() []clients.Tool {
	return []clients.Tool{
		{
			Name:        "list_devices",
			Description: "Returns a list of the user's smart home devices with their IDs, names, rooms, and types. Use this to discover what devices are available before querying state or sending commands.",
			Parameters: `{
				"type": "object",
				"properties": {},
				"required": []
			}`,
		},
		{
			Name:        "get_device_state",
			Description: "Returns the current state of a specific device, including all capabilities (on/off, brightness, color, mode, etc.) and sensor properties (temperature, humidity, etc.).",
			Parameters: `{
				"type": "object",
				"properties": {
					"device_id": {
						"type": "string",
						"description": "The UUID of the device to query"
					}
				},
				"required": ["device_id"]
			}`,
		},
		{
			Name:        "control_device",
			Description: "Sends a control command to a device capability. Use list_devices and get_device_state first to find the device ID and available capability types.",
			Parameters: `{
				"type": "object",
				"properties": {
					"device_id": {
						"type": "string",
						"description": "The UUID of the device to control"
					},
					"capability": {
						"type": "string",
						"description": "The capability type to update, e.g. 'devices.capabilities.on_off', 'devices.capabilities.range', 'devices.capabilities.mode', 'devices.capabilities.toggle', 'devices.capabilities.color_setting'"
					},
					"value": {
						"description": "The new value for the capability. For on_off use true/false. For range use a number. For mode/toggle use a string."
					}
				},
				"required": ["device_id", "capability", "value"]
			}`,
		},
	}
}

func userIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value(ctxkeys.UserID)
	if v == nil {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}
	uid, ok := v.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id has unexpected type in context")
	}
	return uid, nil
}

func listDevicesHandler(repo *deviceRepo) tools.Handler {
	return func(ctx context.Context, _ []byte) (string, error) {
		uid, err := userIDFromCtx(ctx)
		if err != nil {
			return "", err
		}

		devices, err := repo.listDevices(ctx, uid.String())
		if err != nil {
			return "", err
		}

		type item struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Room string `json:"room,omitempty"`
			Type string `json:"type,omitempty"`
		}
		out := make([]item, 0, len(devices))
		for _, d := range devices {
			out = append(out, item{ID: d.ID, Name: d.Name, Room: d.Room, Type: d.Type})
		}

		b, err := json.Marshal(out)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

func getDeviceStateHandler(repo *deviceRepo) tools.Handler {
	return func(ctx context.Context, args []byte) (string, error) {
		uid, err := userIDFromCtx(ctx)
		if err != nil {
			return "", err
		}

		var req struct {
			DeviceID string `json:"device_id"`
		}
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid args: %w", err)
		}
		if req.DeviceID == "" {
			return "", fmt.Errorf("device_id is required")
		}

		device, err := repo.getDevice(ctx, req.DeviceID, uid.String())
		if err != nil {
			return "", err
		}

		type capOut struct {
			Type  string          `json:"type"`
			State json.RawMessage `json:"state,omitempty"`
		}
		type propOut struct {
			Type  string          `json:"type"`
			State json.RawMessage `json:"state,omitempty"`
		}
		type out struct {
			ID           string    `json:"id"`
			Name         string    `json:"name"`
			Room         string    `json:"room,omitempty"`
			Type         string    `json:"type,omitempty"`
			Capabilities []capOut  `json:"capabilities"`
			Properties   []propOut `json:"properties"`
		}

		result := out{
			ID:           device.ID,
			Name:         device.Name,
			Room:         device.Room,
			Type:         device.Type,
			Capabilities: make([]capOut, 0, len(device.Capabilities)),
			Properties:   make([]propOut, 0, len(device.Properties)),
		}
		for _, c := range device.Capabilities {
			result.Capabilities = append(result.Capabilities, capOut{Type: c.Type, State: c.State})
		}
		for _, p := range device.Properties {
			result.Properties = append(result.Properties, propOut{Type: p.Type, State: p.State})
		}

		b, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

func controlDeviceHandler(repo *deviceRepo) tools.Handler {
	return func(ctx context.Context, args []byte) (string, error) {
		uid, err := userIDFromCtx(ctx)
		if err != nil {
			return "", err
		}

		var req struct {
			DeviceID   string          `json:"device_id"`
			Capability string          `json:"capability"`
			Value      json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid args: %w", err)
		}
		if req.DeviceID == "" || req.Capability == "" {
			return "", fmt.Errorf("device_id and capability are required")
		}

		// Verify device belongs to user
		if _, err := repo.getDevice(ctx, req.DeviceID, uid.String()); err != nil {
			return "", fmt.Errorf("device not found: %w", err)
		}

		if err := repo.updateCapabilityState(ctx, req.DeviceID, req.Capability, req.Value); err != nil {
			return "", err
		}

		return fmt.Sprintf(`{"status":"ok","device_id":"%s","capability":"%s"}`, req.DeviceID, req.Capability), nil
	}
}
