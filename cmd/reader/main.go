package main

import (
	"time"

	"github.com/h44z/dntroomloggpro-go/pkg"

	"github.com/sirupsen/logrus"
)

func main() {
	//logrus.SetLevel(logrus.TraceLevel)

	r := pkg.NewRoomLogg(pkg.NewRoomLoggConfig())
	if err := r.Open(); err != nil {
		logrus.Fatal("Unable to initialize DNT RoomLogg PRO!")
	}
	defer r.Close()

	for {
		time.Sleep(30 * time.Second)

		channelData, err1 := r.FetchCurrentData()
		calibrationData, err2 := r.FetchCalibrationData()
		if err1 != nil || err2 != nil {
			logrus.Errorf("Lost connection to DNT RoomLogg PRO: %v, %v", err1, err2)
			r.Close()
			if err := r.Open(); err != nil {
				logrus.Errorf("Failed to restore connection to DNT RoomLogg PRO: %v", err)
			}
			continue
		}

		logrus.Infof("----------------------------------------")
		for _, channel := range channelData {
			logrus.Infof("ChannelData %d:\t %1.1f °C,\t %1.0f %%", channel.Number, channel.Temperature, channel.Humidity)
		}
		for _, channel := range calibrationData {
			logrus.Infof("CalibrationData %d:\t %1.1f °C,\t %1.0f %%", channel.Channel, channel.Temperature, channel.Humidity)
		}
	}
}
