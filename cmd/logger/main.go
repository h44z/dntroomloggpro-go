package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/h44z/dntroomloggpro-go/pkg"

	"github.com/sirupsen/logrus"
)

type publisher interface {
	Publish(settings *pkg.SettingsData, channels []*pkg.ChannelData, isOnline bool) error
}

// You need to create a influx db before using this tool:
/*
$ influx (-username xxx -password yyy)
> create database gmclogg
> exit
*/
func main() {
	logrus.SetLevel(logrus.DebugLevel)

	rCfg := pkg.NewRoomLoggConfig()

	r := pkg.NewRoomLogg(rCfg)
	if err := r.Open(); err != nil {
		logrus.Fatal("[MAIN] Unable to initialize RoomLogg!", err)
	}
	defer r.Close()

	rest, mqtt, influx := features()

	var publishers []publisher

	if rest {
		sCfg := pkg.NewRestConfig()
		s, err := pkg.NewServer(sCfg)
		if err != nil {
			logrus.Fatalf("Unable to initialize WebServer: %v", err)
		}
		s.SetRoomLogInstance(r)
		go s.Run() // start webserver

		publishers = append(publishers, s)
	}

	if mqtt {
		mCfg := pkg.NewMqttConfig()
		p, err := pkg.NewMqttPublisher(mCfg)
		if err != nil {
			logrus.Fatalf("[MAIN] Unable to initialize MQTT publisher: %v", err)
		}
		defer p.Close()

		publishers = append(publishers, p)
	}

	if influx {
		iCfg := pkg.NewInfluxConfig()
		i := pkg.NewInfluxLogger(iCfg)
		defer i.Close()

		publishers = append(publishers, i)
	}

	logrus.Infof("[MAIN] Starting in %v (%d pub)...", time.Duration(rCfg.PollingRate)*time.Second, len(publishers))

	// Start ticker
	ticker := time.NewTicker(time.Duration(rCfg.PollingRate) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			isOnline := true
			settings, err := r.FetchSettings()
			if err != nil {
				logrus.Errorf("[MAIN] Lost connection to RoomLogg: %v", err)
				_ = r.Reconnect()
				isOnline = false
			}
			channelData, err := r.FetchCurrentData()
			if err != nil {
				logrus.Errorf("[MAIN] Lost connection to GCM: %v", err)
				_ = r.Reconnect()
				isOnline = false
			}

			logMsg := make([]string, len(channelData))
			for i, ch := range channelData {
				logMsg[i] = fmt.Sprintf("CH %d: %0.1f°C/%0.0f% ", ch.Number, ch.Temperature, ch.Humidity)
			}
			logrus.Infof("[MAIN] Fetched: %s", strings.Join(logMsg, "; "))

			for i, p := range publishers {
				err := p.Publish(settings, channelData, isOnline)
				if err != nil {
					logrus.Errorf("[MAIN] Failed to publish: %v", err)
				}
				logrus.Debugf("[MAIN] Published #%d", i)
			}

			logrus.Info("[MAIN] Tick completed!")
		}
	}
}

func features() (rest, mqtt, influx bool) {
	rest = true
	mqtt = true
	influx = true

	if val, err := strconv.ParseBool(os.Getenv("ENABLE_REST")); err == nil && !val {
		rest = false
	}
	if val, err := strconv.ParseBool(os.Getenv("ENABLE_MQTT")); err == nil && !val {
		mqtt = false
	}
	if val, err := strconv.ParseBool(os.Getenv("ENABLE_INFLUX")); err == nil && !val {
		influx = false
	}
	return
}
