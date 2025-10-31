package services

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
	"log"
	"strings"
)

// <user-id>/<device-id≥/capabilities/<capability>/set
//<user-id>/<device-id≥/capabilities/<capability>/state
//<user-id>/<device-id>/properties/<property>/state
//<user-id>/<device-id>/state

type deviceListenerService struct {
	mqttClient     mqtt.Client
	devicesService DevicesService
}

func (s *deviceListenerService) StartListener() {
	// Подписываемся на все нужные типы сообщений
	topics := []string{
		"+/+/state",                // общий статус устройства
		"+/+/capabilities/+/state", // состояния "умений"
		"+/+/properties/+/state",   // свойства/сенсоры
	}

	for _, t := range topics {
		if token := s.mqttClient.Subscribe(t, 0, s.handleMessage); token.Wait() && token.Error() != nil {
			log.Printf("cannot subscribe to %s: %v", t, token.Error())
		}
	}
	log.Println("[MQTT] listener started")
}

func (s *deviceListenerService) handleMessage(client mqtt.Client, msg mqtt.Message) {
	_, deviceID, component, name, direction, err := parseTopicExt(msg.Topic())
	if err != nil {
		log.Printf("bad topic: %v", err)
		return
	}

	switch component {
	case "state":
		// Топик <user>/<device>/state
		s.handleDeviceState(deviceID, msg.Payload())
	case "capabilities":
		if direction == "state" {
			s.handleCapabilityState(deviceID, name, msg.Payload())
		}
	case "properties":
		if direction == "state" {
			s.handlePropertyState(deviceID, name, msg.Payload())
		}
	default:
		log.Printf("unknown component in topic: %s", msg.Topic())
	}
}

func (s *deviceListenerService) handleDeviceState(deviceID string, payload []byte) {
	var p struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("bad /state payload for %s: %v", deviceID, err)
		return
	}
	_ = s.devicesService.UpdateLastSeen(deviceID)
}

func (s *deviceListenerService) handleCapabilityState(deviceID, capability string, payload []byte) {
	if err := s.devicesService.UpdateCurrentState(deviceID, capability, payload); err != nil {
		log.Printf("cannot update capability %s for %s: %v", capability, deviceID, err)
	}
}

func (s *deviceListenerService) handlePropertyState(deviceID, property string, payload []byte) {
	if err := s.devicesService.UpdateProperty(deviceID, property, payload); err != nil {
		log.Printf("cannot update property %s for %s: %v", property, deviceID, err)
	}
}

// <user>/<device>/<component>/<name>/<direction?>
func parseTopicExt(topic string) (user, device, component, name, direction string, err error) {
	parts := strings.Split(topic, "/")
	switch len(parts) {
	case 3:
		return parts[0], parts[1], parts[2], "", "", nil
	case 4:
		return parts[0], parts[1], parts[2], parts[3], "", nil
	case 5:
		return parts[0], parts[1], parts[2], parts[3], parts[4], nil
	default:
		return "", "", "", "", "", fmt.Errorf("invalid topic format: %s", topic)
	}
}

func NewStatusListenerService(db *gorm.DB, mqttClient mqtt.Client) StatusListenerService {
	return &deviceListenerService{
		mqttClient:     mqttClient,
		devicesService: NewDevicesService(db),
	}
}
