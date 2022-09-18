package pkg

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/sirupsen/logrus"
)

type MqttPublisher struct {
	// Core components
	cfg    *MqttConfig
	client mqtt.Client
}

func NewMqttPublisher(cfg *MqttConfig) (*MqttPublisher, error) {
	p := &MqttPublisher{
		cfg: cfg,
	}

	err := p.Setup()

	return p, err
}

func (p *MqttPublisher) Setup() error {
	opts := mqtt.NewClientOptions()
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(2 * time.Second)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", p.cfg.Broker, p.cfg.Port))
	opts.SetClientID(fmt.Sprintf("roomlogg_mqtt_%s", p.cfg.Topic))
	if p.cfg.Username != "" {
		opts.SetUsername(p.cfg.Username)
	}
	if p.cfg.Password != "" {
		opts.SetPassword(p.cfg.Password)
	}
	opts.SetDefaultPublishHandler(p.onMessageReceived)
	opts.OnConnect = p.onConnectHandler
	opts.OnConnectionLost = p.onConnectionLostHandler
	p.client = mqtt.NewClient(opts)

	if token := p.client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("[MQTT] Setup of mqtt publisher failed: %v!", token.Error())
		panic(token.Error())
	}

	logrus.Infof("[MQTT] Setup of mqtt publisher completed!")
	return nil
}

func (p *MqttPublisher) Close() {
	p.client.Disconnect(250)
}

func (p *MqttPublisher) onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	logrus.Infof("[MQTT] TOPIC: %s", msg.Topic())
	logrus.Infof("[MQTT] MSG: %s", msg.Payload())
}

func (p *MqttPublisher) onConnectHandler(_ mqtt.Client) {
	logrus.Infof("[MQTT] Connected to broker!")
}

func (p *MqttPublisher) onConnectionLostHandler(_ mqtt.Client, err error) {
	logrus.Warnf("[MQTT] Connection to broker lost: %v!", err)
}

func (p *MqttPublisher) Publish(settings *SettingsData, channels []*ChannelData, isOnline bool) error {
	if err := p.publishHomeAssistantConfig(channels); err != nil {
		return fmt.Errorf("failed to publish mqtt config: %w", err)
	}

	time.Sleep(2 * time.Second) // wait for home assistant to process new topics

	if err := p.publishTopics(channels, isOnline); err != nil {
		return fmt.Errorf("failed to publish mqtt sensors: %w", err)
	}

	return nil
}

func (p *MqttPublisher) publishHomeAssistantConfig(channels []*ChannelData) error {
	topicStatus := fmt.Sprintf("homeassistant/binary_sensor/%s/status/config", p.cfg.Topic)
	availabilityConfig := map[string]any{
		"name":               "Status",
		"state_topic":        fmt.Sprintf("roomlogg/%s/status", p.cfg.Topic),
		"availability_topic": fmt.Sprintf("roomlogg/%s/status", p.cfg.Topic),
		"device_class":       "connectivity",
		"payload_on":         "online",
		"payload_off":        "offline",
		"expire_after":       "240",
		"unique_id":          fmt.Sprintf("roomlogg_%s_status", p.cfg.Topic),
		"device": map[string]any{
			"identifiers":  p.cfg.Topic,
			"name":         p.cfg.Topic,
			"manufacturer": "DNT",
			"model":        "DNT RoomLogg PRO",
		},
	}

	payload, _ := json.Marshal(availabilityConfig)
	token := p.client.Publish(topicStatus, 0, false, string(payload))
	token.Wait()

	for _, ch := range channels {
		topicTemperature := fmt.Sprintf("homeassistant/sensor/%s/temperature_%d/config", p.cfg.Topic, ch.Number)
		temperatureConfig := map[string]any{
			"name":                fmt.Sprintf("Temperature Channel %d", ch.Number),
			"state_topic":         fmt.Sprintf("roomlogg/%s/temperature/%d", p.cfg.Topic, ch.Number),
			"availability_topic":  fmt.Sprintf("roomlogg/%s/status", p.cfg.Topic),
			"unit_of_measurement": "°C",
			"device_class":        "temperature",
			"state_class":         "measurement",
			"value_template":      "{{ value_json.value | float }}",
			"unique_id":           fmt.Sprintf("roomlogg_%s_temp_%d", p.cfg.Topic, ch.Number),
			"device": map[string]any{
				"identifiers":  p.cfg.Topic,
				"name":         p.cfg.Topic,
				"manufacturer": "DNT",
				"model":        "DNT RoomLogg PRO",
			},
		}
		payload, _ = json.Marshal(temperatureConfig)
		token = p.client.Publish(topicTemperature, 0, false, string(payload))
		token.Wait()

		topicHumidity := fmt.Sprintf("homeassistant/sensor/%s/humidity_%d/config", p.cfg.Topic, ch.Number)
		humidityConfig := map[string]any{
			"name":                fmt.Sprintf("Humidity Channel %d", ch.Number),
			"state_topic":         fmt.Sprintf("roomlogg/%s/humidity/%d", p.cfg.Topic, ch.Number),
			"availability_topic":  fmt.Sprintf("roomlogg/%s/status", p.cfg.Topic),
			"unit_of_measurement": "%",
			"device_class":        "humidity",
			"state_class":         "measurement",
			"value_template":      "{{ value_json.value | float }}",
			"unique_id":           fmt.Sprintf("roomlogg_%s_humid_%d", p.cfg.Topic, ch.Number),
			"device": map[string]any{
				"identifiers":  p.cfg.Topic,
				"name":         p.cfg.Topic,
				"manufacturer": "DNT",
				"model":        "DNT RoomLogg PRO",
			},
		}
		payload, _ = json.Marshal(humidityConfig)
		token = p.client.Publish(topicHumidity, 0, false, string(payload))
		token.Wait()
	}
	return nil
}

func (p *MqttPublisher) publishTopics(channels []*ChannelData, isOnline bool) error {
	topicStatus := fmt.Sprintf("roomlogg/%s/status", p.cfg.Topic)
	status := "offline"
	if isOnline {
		status = "online"
	}
	token := p.client.Publish(topicStatus, 0, false, status)
	token.Wait()

	for _, ch := range channels {
		topicTemperature := fmt.Sprintf("roomlogg/%s/temperature/%d", p.cfg.Topic, ch.Number)
		temperatureValue := map[string]any{
			"value":   ch.Temperature,
			"unit":    "°C",
			"channel": ch.Number,
		}
		payload, _ := json.Marshal(temperatureValue)
		token = p.client.Publish(topicTemperature, 0, false, string(payload))
		token.Wait()

		topicHumidity := fmt.Sprintf("roomlogg/%s/humidity/%d", p.cfg.Topic, ch.Number)
		humidityValue := map[string]any{
			"value":   ch.Humidity,
			"unit":    "%",
			"channel": ch.Number,
		}
		payload, _ = json.Marshal(humidityValue)
		token = p.client.Publish(topicHumidity, 0, false, string(payload))
		token.Wait()
	}
	return nil
}
