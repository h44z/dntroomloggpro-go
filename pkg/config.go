package pkg

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type InfluxConfig struct {
	URL      string `envconfig:"INFLUX_URL"`
	UserName string `envconfig:"INFLUX_USER"`
	Password string `envconfig:"INFLUX_PASS"`
	Bucket   string `envconfig:"INFLUX_BUCKET"`
}

func NewInfluxConfig() *InfluxConfig {
	cfg := &InfluxConfig{
		URL:      "http://localhost:8086",
		UserName: "influxuser",
		Password: "influxpass",
		Bucket:   "",
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
