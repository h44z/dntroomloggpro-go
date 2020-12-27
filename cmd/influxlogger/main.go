package main

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/h44z/dntroomloggpro-go/pkg"
	"github.com/sirupsen/logrus"
)

// You need to create a influx db before using this tool:
/*
$ influx (-username xxx -password yyy)
> create database roomlogg
> exit
*/
func main() {
	//logrus.SetLevel(logrus.TraceLevel)

	r := pkg.NewRoomLogg()
	if err := r.Open(); err != nil {
		logrus.Fatal("Unable to initialize DNT RoomLogg PRO!")
	}
	defer r.Close()

	iCfg := pkg.NewInfluxConfig()
	i := pkg.NewInfluxLogger(iCfg)
	defer i.Close()

	for {
		time.Sleep(30 * time.Second)

		settings, err1 := r.FetchSettings()
		channelData, err2 := r.FetchCurrentData()
		if err1 != nil || err2 != nil {
			logrus.Errorf("Lost connection to DNT RoomLogg PRO: %v, %v", err1, err2)
			r.Close()
			if err := r.Open(); err != nil {
				logrus.Errorf("Failed to restore connection to DNT RoomLogg PRO: %v", err)
			}
			continue
		}

		tempUnit := "°C"
		if settings.Units == pkg.UnitFahrenheit {
			tempUnit = "°F"
		}

		points := make([]*write.Point, 0, len(channelData)*2)
		for _, channel := range channelData {
			points = append(points, influxdb2.NewPoint("temperature", // Measurement
				map[string]string{"unit": tempUnit, "channel": fmt.Sprintf("%d", channel.Number)}, // Tags
				map[string]interface{}{"value": channel.Temperature},                              // Fields
				time.Now()))
			points = append(points, influxdb2.NewPoint("humidity", // Measurement
				map[string]string{"unit": "%", "channel": fmt.Sprintf("%d", channel.Number)}, // Tags
				map[string]interface{}{"value": channel.Humidity},                            // Fields
				time.Now()))
		}
		if err := i.LogPoints(iCfg.Bucket, points...); err != nil {
			logrus.Errorf("Lost connection to InfluxDB: %v", err)
		}
	}
}
