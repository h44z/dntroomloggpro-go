package main

import (
	"github.com/h44z/dntroomloggpro-go/pkg"
	"github.com/sirupsen/logrus"
)

func main() {
	//logrus.SetLevel(logrus.TraceLevel)

	sCfg := pkg.NewServerConfig()
	s, err := pkg.NewServer(sCfg)
	if err != nil {
		logrus.Fatalf("Unable to initialize WebServer: %v", err)
	}
	defer s.Close()

	s.Run()
}
