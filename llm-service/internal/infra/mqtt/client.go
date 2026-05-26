package mqtt

import (
	"fmt"
	"llm-service/internal/utils"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	log    *slog.Logger
	Client mqtt.Client
}

func NewClient(log *slog.Logger) (*Client, error) {
	brokerURL := utils.GetEnv("MQTT_URL", "tcp://broker:1883")
	user := utils.GetEnv("MQTT_USER", "user")
	pass := utils.GetEnv("MQTT_PASSWORD", "superdifficultpassword")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(fmt.Sprintf("llm-service-%d", time.Now().UnixNano()))
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.OnConnect = func(_ mqtt.Client) {
		log.Info("mqtt connected", slog.String("broker", brokerURL))
	}
	opts.OnConnectionLost = func(_ mqtt.Client, err error) {
		log.Warn("mqtt connection lost", slog.Any("err", err))
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("mqtt connect timeout")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("mqtt connect: %w", err)
	}

	return &Client{log: log, Client: client}, nil
}

func (c *Client) Close() {
	if c.Client != nil {
		c.Client.Disconnect(250)
	}
}
