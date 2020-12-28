package pkg

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type InfluxConfig struct {
	URL             string `envconfig:"INFLUX_URL"`
	UserName        string `envconfig:"INFLUX_USER"`
	Password        string `envconfig:"INFLUX_PASS"`
	Bucket          string `envconfig:"INFLUX_BUCKET"`
	IntervalSeconds int    `envconfig:"INFLUX_INTERVAL"`
}

func NewInfluxConfig() *InfluxConfig {
	// Default config
	cfg := &InfluxConfig{
		URL:             "http://localhost:8086",
		UserName:        "influxuser",
		Password:        "influxpass",
		Bucket:          "roomlogg",
		IntervalSeconds: 60,
	}
	if err := loadConfigEnv(cfg); err != nil {
		logrus.Warnf("unable to load environment config: %v", err)
	}

	return cfg
}

type ServerConfig struct {
	ListenAddress string `envconfig:"RESTAPI_ADDRESS"`
}

func NewServerConfig() *ServerConfig {
	// Default config
	cfg := &ServerConfig{
		ListenAddress: ":8080",
	}
	if err := loadConfigEnv(cfg); err != nil {
		logrus.Warnf("unable to load environment config: %v", err)
	}

	return cfg
}

func loadConfigEnv(cfg interface{}) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}

	return nil
}
