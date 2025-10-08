package storage

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnectMQTT(opts *mqtt.ClientOptions) (mqtt.Client, error) {
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
