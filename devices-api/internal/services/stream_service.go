package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceEvent struct {
	Type      string          `json:"type"`
	DeviceID  string          `json:"device_id"`
	Capability string         `json:"capability,omitempty"`
	Property   string         `json:"property,omitempty"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"ts"`
}

type StreamService interface {
	Subscribe(ctx context.Context, userID uuid.UUID) (<-chan DeviceEvent, func())
}

type streamService struct {
	mqtt mqtt.Client
	log  *slog.Logger

	mu          sync.Mutex
	subscribers map[uuid.UUID]map[int64]chan DeviceEvent
	nextID      int64
	mqttSubs    map[uuid.UUID]int // ref-count per user (number of subscribers)
}

func NewStreamService(mqttClient mqtt.Client, log *slog.Logger) StreamService {
	return &streamService{
		mqtt:        mqttClient,
		log:         log,
		subscribers: make(map[uuid.UUID]map[int64]chan DeviceEvent),
		mqttSubs:    make(map[uuid.UUID]int),
	}
}

func (s *streamService) Subscribe(ctx context.Context, userID uuid.UUID) (<-chan DeviceEvent, func()) {
	ch := make(chan DeviceEvent, 32)

	s.mu.Lock()
	s.nextID++
	id := s.nextID
	if _, ok := s.subscribers[userID]; !ok {
		s.subscribers[userID] = make(map[int64]chan DeviceEvent)
	}
	s.subscribers[userID][id] = ch
	firstForUser := s.mqttSubs[userID] == 0
	s.mqttSubs[userID]++
	s.mu.Unlock()

	if firstForUser {
		s.subscribeMQTT(userID)
	}

	cancel := func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		if subs, ok := s.subscribers[userID]; ok {
			if c, ok := subs[id]; ok {
				close(c)
				delete(subs, id)
			}
			if len(subs) == 0 {
				delete(s.subscribers, userID)
			}
		}

		s.mqttSubs[userID]--
		if s.mqttSubs[userID] <= 0 {
			delete(s.mqttSubs, userID)
			s.unsubscribeMQTT(userID)
		}
	}

	go func() {
		<-ctx.Done()
		cancel()
	}()

	return ch, cancel
}

func (s *streamService) subscribeMQTT(userID uuid.UUID) {
	topic := fmt.Sprintf("%s/#", userID.String())
	s.log.Info("mqtt subscribe", slog.String("topic", topic))

	token := s.mqtt.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
		s.handleMessage(userID, msg.Topic(), msg.Payload())
	})
	if !token.WaitTimeout(5 * time.Second) {
		s.log.Warn("mqtt subscribe timeout", slog.String("topic", topic))
		return
	}
	if err := token.Error(); err != nil {
		s.log.Warn("mqtt subscribe failed", slog.String("topic", topic), slog.Any("err", err))
	}
}

func (s *streamService) unsubscribeMQTT(userID uuid.UUID) {
	topic := fmt.Sprintf("%s/#", userID.String())
	s.log.Info("mqtt unsubscribe", slog.String("topic", topic))
	token := s.mqtt.Unsubscribe(topic)
	token.WaitTimeout(2 * time.Second)
}

// handleMessage parses topic into event and forwards to subscribers
// Topic layouts handled:
//   <user_id>/<device_id>/capabilities/<type>/state
//   <user_id>/<device_id>/properties/<type>/state
//   <user_id>/<device_id>/state
//   <user_id>/<device_id>/status
func (s *streamService) handleMessage(userID uuid.UUID, topic string, payload []byte) {
	parts := strings.Split(topic, "/")
	if len(parts) < 3 {
		return
	}
	deviceID := parts[1]

	var ev DeviceEvent
	ev.DeviceID = deviceID
	ev.Timestamp = time.Now().UnixMilli()
	ev.Payload = json.RawMessage(payload)

	switch {
	case len(parts) >= 5 && parts[2] == "capabilities" && parts[4] == "state":
		ev.Type = "capability_state"
		ev.Capability = parts[3]
	case len(parts) >= 5 && parts[2] == "properties" && parts[4] == "state":
		ev.Type = "property_state"
		ev.Property = parts[3]
	case len(parts) == 3 && parts[2] == "state":
		ev.Type = "device_state"
	case len(parts) == 3 && parts[2] == "status":
		ev.Type = "device_status"
	default:
		// Ignore set commands and other topics we publish ourselves
		return
	}

	s.mu.Lock()
	subs, ok := s.subscribers[userID]
	if !ok {
		s.mu.Unlock()
		return
	}
	// Snapshot the channels to release the lock quickly
	chans := make([]chan DeviceEvent, 0, len(subs))
	for _, c := range subs {
		chans = append(chans, c)
	}
	s.mu.Unlock()

	for _, c := range chans {
		select {
		case c <- ev:
		default:
			// Subscriber is too slow; drop to avoid blocking the MQTT thread
		}
	}
}
