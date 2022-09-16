package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/sirupsen/logrus"
)

type MqttPublisher struct {
	// Core components
	config  *MqttConfig
	client  mqtt.Client
	station *RoomLogg

	configInterval time.Duration
	sensorInterval time.Duration
}

func NewMqttPublisher(cfg *MqttConfig) (*MqttPublisher, error) {
	p := &MqttPublisher{}

	// Setup base station
	p.station = NewRoomLogg()
	if p.station == nil {
		return nil, errors.New("failed to initialize RoomLogg PRO")
	}
	if err := p.station.Open(); err != nil {
		return nil, errors.New("failed to initialize RoomLogg PRO, open failed")
	}

	err := p.Setup(cfg)

	return p, err
}

func NewMqttPublisherWithBaseStation(cfg *MqttConfig, station *RoomLogg) (*MqttPublisher, error) {
	p := &MqttPublisher{}

	// Setup base station
	p.station = station

	err := p.Setup(cfg)

	return p, err
}

func (p *MqttPublisher) Setup(cfg *MqttConfig) error {
	p.config = cfg
	p.sensorInterval = 1 * time.Minute
	p.configInterval = 10 * time.Minute

	opts := mqtt.NewClientOptions()
	opts.SetKeepAlive(60 * time.Second)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Broker, cfg.Port))
	opts.SetClientID("roomlogg_mqtt")
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	opts.OnConnect = p.onConnectHandler
	opts.OnConnectionLost = p.onConnectionLostHandler
	p.client = mqtt.NewClient(opts)

	logrus.Infof("Setup of mqtt publisher completed!")
	return nil
}

func (p *MqttPublisher) Close() {
	p.station.Close()
	p.client.Disconnect(250)
}

func (p *MqttPublisher) onConnectHandler(_ mqtt.Client) {
	logrus.Infof("[MQTT] Connected to broker!")
}

func (p *MqttPublisher) onConnectionLostHandler(_ mqtt.Client, err error) {
	logrus.Warnf("[MQTT] Connection to broker lost: %v!", err)
}

func (p *MqttPublisher) Run() {

	// Initial tick
	if err := p.configTick(); err != nil {
		logrus.Errorf("[MQTT] initial config tick failed: %v", err)
	}

	// Start ticker
	tickerSensor := time.NewTicker(p.sensorInterval)
	tickerConfig := time.NewTicker(p.configInterval)
	defer tickerSensor.Stop()
	defer tickerConfig.Stop()
	for {
		select {
		case <-tickerConfig.C:
			if err := p.configTick(); err != nil {
				logrus.Errorf("[MQTT] config tick failed: %v", err)
			}
		case <-tickerSensor.C:
			if err := p.sensorTick(); err != nil {
				logrus.Errorf("[MQTT] sensor tick failed: %v", err)
			}
		}
	}
}

func (p *MqttPublisher) sensorTick() error {
	online := true
	data, err := p.station.FetchCurrentData()
	if err != nil {
		online = false
		logrus.Errorf("[MQTT] Lost connection to DNT RoomLogg PRO: %v", err)
		p.station.Close()
		if err := p.station.Open(); err != nil {
			logrus.Errorf("[MQTT] Failed to restore connection to DNT RoomLogg PRO: %v", err)
		}
	}

	if err := p.publishTopics(data, online); err != nil {
		return err
	}
	return nil
}

func (p *MqttPublisher) configTick() error {
	data, err := p.station.FetchCurrentData()
	if err != nil {
		logrus.Errorf("[MQTT] Lost connection to DNT RoomLogg PRO: %v", err)
		p.station.Close()
		if err := p.station.Open(); err != nil {
			logrus.Errorf("[MQTT] Failed to restore connection to DNT RoomLogg PRO: %v", err)
			return err
		}

		return errors.New("reconnect skip")
	}

	if err := p.publishHomeAssistantConfig(data); err != nil {
		return err
	}
	return nil
}

func (p *MqttPublisher) publishHomeAssistantConfig(channels []*ChannelData) error {
	topicStatus := fmt.Sprintf("homeassistant/binary_sensor/%s/status/config", p.config.Topic)
	availabilityConfig := map[string]any{
		"name":               "Status",
		"state_topic":        fmt.Sprintf("roomlogg/%s/status", p.config.Topic),
		"availability_topic": fmt.Sprintf("roomlogg/%s/status", p.config.Topic),
		"device_class":       "connectivity",
		"payload_on":         "online",
		"payload_off":        "offline",
		"expire_after":       "120",
		"unique_id":          fmt.Sprintf("roomlogg_%s_status", p.config.Topic),
		"device": map[string]any{
			"identifiers":  p.config.Topic,
			"name":         p.config.Topic,
			"manufacturer": "DNT",
			"model":        "DNT RoomLogg PRO",
		},
	}

	payload, _ := json.Marshal(availabilityConfig)
	token := p.client.Publish(topicStatus, 0, false, string(payload))
	token.Wait()

	for _, ch := range channels {
		topicTemperature := fmt.Sprintf("homeassistant/sensor/%s/temperature_%d/config", p.config.Topic, ch.Number)
		temperatureConfig := map[string]any{
			"name":                fmt.Sprintf("Temperature Channel %d", ch.Number),
			"state_topic":         fmt.Sprintf("roomlogg/%s/temperature/%d", p.config.Topic, ch.Number),
			"availability_topic":  fmt.Sprintf("roomlogg/%s/status", p.config.Topic),
			"unit_of_measurement": "°C",
			"device_class":        "temperature",
			"state_class":         "measurement",
			"value_template":      "{{ value_json.value | float }}",
			"unique_id":           fmt.Sprintf("roomlogg_%s_temp_%d", p.config.Topic, ch.Number),
			"device": map[string]any{
				"identifiers":  p.config.Topic,
				"name":         p.config.Topic,
				"manufacturer": "DNT",
				"model":        "DNT RoomLogg PRO",
			},
		}
		payload, _ = json.Marshal(temperatureConfig)
		token = p.client.Publish(topicTemperature, 0, false, string(payload))
		token.Wait()

		topicHumidity := fmt.Sprintf("homeassistant/sensor/%s/humidity_%d/config", p.config.Topic, ch.Number)
		humidityConfig := map[string]any{
			"name":                fmt.Sprintf("Temperature Channel %d", ch.Number),
			"state_topic":         fmt.Sprintf("roomlogg/%s/humidity/%d", p.config.Topic, ch.Number),
			"availability_topic":  fmt.Sprintf("roomlogg/%s/status", p.config.Topic),
			"unit_of_measurement": "%",
			"device_class":        "humidity",
			"state_class":         "measurement",
			"value_template":      "{{ value_json.value | float }}",
			"unique_id":           fmt.Sprintf("roomlogg_%s_humid_%d", p.config.Topic, ch.Number),
			"device": map[string]any{
				"identifiers":  p.config.Topic,
				"name":         p.config.Topic,
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
	topicStatus := fmt.Sprintf("roomlogg/%s/status", p.config.Topic)
	status := "offline"
	if isOnline {
		status = "online"
	}
	token := p.client.Publish(topicStatus, 0, false, status)
	token.Wait()

	for _, ch := range channels {
		topicTemperature := fmt.Sprintf("roomlogg/%s/temperature/%d", p.config.Topic, ch.Number)
		temperatureValue := map[string]any{
			"value":   ch.Temperature,
			"unit":    "°C",
			"channel": ch.Number,
		}
		payload, _ := json.Marshal(temperatureValue)
		token = p.client.Publish(topicTemperature, 0, false, string(payload))
		token.Wait()

		topicHumidity := fmt.Sprintf("roomlogg/%s/humidity/%d", p.config.Topic, ch.Number)
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
