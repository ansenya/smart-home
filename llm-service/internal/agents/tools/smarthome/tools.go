package smarthome

import (
	"context"
	"encoding/json"
	"fmt"
	"llm-service/internal/agents/tools"
	"llm-service/internal/clients"
	"llm-service/internal/ctxkeys"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, mqttClient mqtt.Client, reg tools.Registry) []clients.Tool {
	repo := newDeviceRepo(db, mqttClient)

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
			Description: "Sends a control command to a device capability. The command is delivered to the physical device over MQTT. Use list_devices and get_device_state first to find the device ID, available capability types, and the correct instance for that capability. For on_off capability the instance is always 'on'. For color_setting use 'hsv' (value: {h:0-360, s:0-100, v:0-100}), 'temperature_k' (value: 2700-6500), or 'scene' (value: scene id). For range/mode/toggle capabilities the instance depends on the device — read it from get_device_state first.",
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
					"instance": {
						"type": "string",
						"description": "Optional capability instance (e.g. 'on', 'brightness', 'hsv', 'temperature_k', 'scene'). Defaults to 'on' for on_off and to the current instance from get_device_state for other capabilities."
					},
					"value": {
						"description": "The new value for the capability. For on_off use true/false. For range use a number. For mode/toggle use a string. For color_setting hsv use {h,s,v}; for temperature_k use a number; for scene use a string."
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

func defaultInstance(capType string) string {
	switch capType {
	case "devices.capabilities.on_off":
		return "on"
	case "devices.capabilities.color_setting":
		return "hsv"
	default:
		return ""
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
			Instance   string          `json:"instance"`
			Value      json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal(args, &req); err != nil {
			return "", fmt.Errorf("invalid args: %w", err)
		}
		if req.DeviceID == "" || req.Capability == "" {
			return "", fmt.Errorf("device_id and capability are required")
		}

		if _, err := repo.getDevice(ctx, req.DeviceID, uid.String()); err != nil {
			return "", fmt.Errorf("device not found: %w", err)
		}

		instance := req.Instance
		if instance == "" {
			if existing, ok := repo.readInstance(ctx, req.DeviceID, req.Capability); ok {
				instance = existing
			} else {
				instance = defaultInstance(req.Capability)
			}
		}

		if err := repo.updateCapabilityState(ctx, req.DeviceID, req.Capability, instance, req.Value); err != nil {
			return "", err
		}

		if err := repo.publishSet(uid.String(), req.DeviceID, req.Capability, instance, req.Value); err != nil {
			return "", fmt.Errorf("mqtt publish: %w", err)
		}

		return fmt.Sprintf(`{"status":"ok","device_id":"%s","capability":"%s","instance":"%s"}`, req.DeviceID, req.Capability, instance), nil
	}
}
