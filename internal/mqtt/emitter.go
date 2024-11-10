package mqtt

import (
	"context"
	"fmt"

	"github.com/denwilliams/go-wiim-mqtt/internal/logging"
)

func NewMqttStatusEmitter(client *MQTTClient) *MqttStatusEmitter {
	return &MqttStatusEmitter{client: client}
}

type MqttStatusEmitter struct {
	client *MQTTClient
}

func (e *MqttStatusEmitter) EmitStatus(ctx context.Context, id string, statusKey string, data interface{}) error {
	topic := fmt.Sprintf("/status/%s/%s", id, statusKey)
	logging.Info("Publishing to %s", topic)
	return e.client.Publish(topic, data)
}
