package pkg

import (
	"context"
	"fmt"

	"github.com/influxdata/influxdb-client-go/v2/api/write"

	infuxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxLogger struct {
	client infuxdb2.Client
}

func NewInfluxLogger(cfg *InfluxConfig) *InfluxLogger {
	i := &InfluxLogger{}
	i.client = infuxdb2.NewClient(cfg.URL, fmt.Sprintf("%s:%s", cfg.UserName, cfg.Password))
	return i
}

func (l *InfluxLogger) Close() {
	if l.client != nil {
		l.client.Close()
	}
}

func (l *InfluxLogger) LogPoints(bucket string, points ...*write.Point) error {
	writeAPI := l.client.WriteAPIBlocking("", bucket)

	// Write data
	err := writeAPI.WritePoint(context.Background(), points...)
	if err != nil {
		return err
	}

	return nil
}

func (l *InfluxLogger) LogCurrentData() error {

	return nil
}
