package mqtt

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/denwilliams/go-wiim-mqtt/internal/logging"
	pm "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(uri *url.URL, baseTopic string, subscribeTopic string) *MQTTClient {
	opts := pm.NewClientOptions().AddBroker(uri.String()).SetClientID("wiim_mqtt_" + uniuri.New()).SetOnConnectHandler(onConnectHandler).SetConnectionLostHandler(onConnectionLostHandler)

	client := pm.NewClient(opts)
	return &MQTTClient{client: &client, baseTopic: baseTopic, subscribeTopic: subscribeTopic}
}

type MQTTClient struct {
	client         *pm.Client
	baseTopic      string
	subscribeTopic string
}

func (mc *MQTTClient) Publish(topic string, data interface{}) error {
	payload, err := serializePayload(data)
	if err != nil {
		return err
	}

	fullTopic := mc.baseTopic + topic

	// Publish a message to the topic with a QoS of 1
	if token := (*mc.client).Publish(fullTopic, 1, false, payload); token.Wait() && token.Error() != nil {
		// TODO: don't panic, just return
		// panic(token.Error())
		logging.Warn("Error publishing message: %s", token.Error())
		return err
	}

	return nil
}

func (mc *MQTTClient) Connect(h CommandHandler) {
	// Connect to the MQTT broker
	if token := (*mc.client).Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// logging.Info("Connected to MQTT")

	prefix := strings.Replace(mc.subscribeTopic, "#", "", 1)

	// Set up a callback function to handle incoming messages
	messageHandler := func(client pm.Client, msg pm.Message) {
		topic := msg.Topic()
		if !strings.HasPrefix(topic, prefix) {
			return
		}
		path := strings.Replace(topic, prefix, "", 1)
		parts := strings.Split(path, "/")
		l := len(parts)

		if l < 2 {
			return
		}
		name := parts[0]
		cmd := parts[1]
		var arg *string = nil
		if l > 2 {
			arg = &parts[2]
		}

		bytes := msg.Payload()
		// payload, err := parsePayload(&bytes)
		// if err != nil {
		// 	logging.Warn("Error unmarshalling JSON: %s %v", err, string(bytes))
		// 	return
		// }
		// logging.Debug("Received message on topic %s: %s", id, payload.String())

		// messages := make(chan string)
		// go func() { messages <- "ping" }()
		// msg := <-messages

		go func() {
			h.HandleCommand(name, cmd, arg, nil, &bytes)
		}()
	}

	// Subscribe to the topic with a QoS of 1
	if token := (*mc.client).Subscribe(mc.subscribeTopic, 1, messageHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logging.Info("Subscribed to %s", mc.subscribeTopic)
}

func (mc *MQTTClient) Disconnect() {
	logging.Info("Disconnecting from MQTT")

	// Unsubscribe from the topic
	if token := (*mc.client).Unsubscribe(mc.baseTopic); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Disconnect from the MQTT broker
	(*mc.client).Disconnect(250)
}

// func parsePayload(bytes *[]byte) (*Command, error) {
// 	var payload Command
// 	if err := json.Unmarshal(*bytes, &payload); err == nil {
// 		return &payload, nil
// 	}

// 	var passOne string
// 	if err := json.Unmarshal(*bytes, &passOne); err != nil {
// 		return nil, err
// 	}

// 	if err := json.Unmarshal([]byte(passOne), &payload); err != nil {
// 		return nil, err
// 	}

// 	return &payload, nil
// }

func serializePayload(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

func onConnectHandler(c pm.Client) {
	logging.Info("Connected to MQTT")
}

func onConnectionLostHandler(c pm.Client, err error) {
	panic(err.Error())
}
