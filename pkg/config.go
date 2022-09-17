package pkg

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type RoomLoggConfig struct {
	PollingRate int `envconfig:"POLLING_RATE"` // Seconds
}

func NewRoomLoggConfig() *RoomLoggConfig {
	// Default config
	cfg := &RoomLoggConfig{
		PollingRate: 60, // 1 Minute
	}
	if err := loadConfigEnv(cfg); err != nil {
		logrus.Warnf("unable to load environment config: %v", err)
	}

	return cfg
}

type InfluxConfig struct {
	URL      string `envconfig:"INFLUX_URL"`
	UserName string `envconfig:"INFLUX_USER"`
	Password string `envconfig:"INFLUX_PASS"`
	Bucket   string `envconfig:"INFLUX_BUCKET"`
}

func NewInfluxConfig() *InfluxConfig {
	// Default config
	cfg := &InfluxConfig{
		URL:      "http://localhost:8086",
		UserName: "influxuser",
		Password: "influxpass",
		Bucket:   "roomlogg",
	}
	if err := loadConfigEnv(cfg); err != nil {
		logrus.Warnf("unable to load environment config: %v", err)
	}

	return cfg
}

type RestConfig struct {
	ListenAddress string `envconfig:"RESTAPI_ADDRESS"`
}

func NewRestConfig() *RestConfig {
	// Default config
	cfg := &RestConfig{
		ListenAddress: ":8080",
	}
	if err := loadConfigEnv(cfg); err != nil {
		logrus.Warnf("unable to load environment config: %v", err)
	}

	return cfg
}

type MqttConfig struct {
	Broker   string `envconfig:"MQTT_BROKER"`
	Port     int    `envconfig:"MQTT_PORT"`
	Username string `envconfig:"MQTT_USER"`
	Password string `envconfig:"MQTT_PASS"`

	Topic string `envconfig:"MQTT_TOPIC"`
}

func NewMqttConfig() *MqttConfig {
	// Default config
	cfg := &MqttConfig{
		Broker:   "localhost",
		Port:     1883,
		Username: "mqttUser",
		Password: "mqttPassword",
		Topic:    "roomlogg",
	}
	if err := loadConfigEnv(cfg); err != nil {
		logrus.Warnf("unable to load environment config: %v", err)
	}

	return cfg
}

func loadConfigEnv(cfg any) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}

	return nil
}
