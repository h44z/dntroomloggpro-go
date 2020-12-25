package main

import (
	"time"

	"github.com/h44z/dntroomloggpro-go/internal"

	"github.com/sirupsen/logrus"
)

func main() {
	//logrus.SetLevel(logrus.TraceLevel)

	r := internal.NewRoomlog()
	if err := r.Open(); err != nil {
		logrus.Fatal("Unable to initialize DNT RoomLogg PRO!")
	}
	defer r.Close()

	for {
		time.Sleep(30 * time.Second)

		channelData, err := r.FetchData()
		if err != nil {
			logrus.Errorf("Lost connection to DNT RoomLogg PRO: %v", err)
			r.Close()
			if err := r.Open(); err != nil {
				logrus.Errorf("Failed to restore connection to DNT RoomLogg PRO: %v", err)
			}
			continue
		}

		logrus.Infof("----------------------------------------")
		for _, channel := range channelData {
			logrus.Infof("Channel %d:\t %1.1f °C,\t %1.0f %%", channel.Number, channel.Temperature, channel.Humidity)
		}
	}
}
