package pkg

import (
	"errors"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Server struct {
	// Core components
	config  *ServerConfig
	server  *gin.Engine
	station *RoomLogg
}

func getExecutableDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Errorf("Failed to get executable directory: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "assets")); os.IsNotExist(err) {
		return "." // assets directory not found -> we are developing in goland =)
	}

	return dir
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	s := &Server{}

	// Setup base station
	s.station = NewRoomLogg()
	if s.station == nil {
		return nil, errors.New("failed to initialize RoomLogg PRO")
	}
	if err := s.station.Open(); err != nil {
		return nil, errors.New("failed to initialize RoomLogg PRO, open failed")
	}

	err := s.Setup(cfg)

	return s, err
}

func NewServerWithBaseStation(cfg *ServerConfig, station *RoomLogg) (*Server, error) {
	s := &Server{}

	// Setup base station
	s.station = station

	err := s.Setup(cfg)

	return s, err
}

func (s *Server) Setup(cfg *ServerConfig) error {
	dir := getExecutableDirectory()
	rDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	logrus.Infof("Real working directory: %s", rDir)
	logrus.Infof("Current working directory: %s", dir)

	// Init rand
	rand.Seed(time.Now().UnixNano())

	// Setup http server
	s.server = gin.Default()
	s.config = cfg

	// Setup all routes
	s.server.GET("/", s.GetCurrentData)
	s.server.GET("/calibration", s.GetCalibrationData)
	s.server.POST("/calibration", s.SetCalibration)
	s.server.GET("/settings", s.GetSettings)
	s.server.POST("/settings", s.SetSettings)
	s.server.GET("/alarm-settings", s.GetAlarmSettings)
	s.server.POST("/alarm-settings", s.SetAlarmSettings)
	s.server.GET("/temperature-alarms", s.GetTemperatureAlarms)
	s.server.POST("/temperature-alarms", s.SetTemperatureAlarms)
	s.server.GET("/humidity-alarms", s.GetHumidityAlarms)
	s.server.POST("/humidity-alarms", s.SetHumidityAlarms)
	s.server.GET("/interval", s.GetIntervalMinutes)
	s.server.POST("/interval", s.SetIntervalMinutes)
	s.server.POST("/language", s.SetLanguage)
	s.server.POST("/time", s.SetCurrentTime)

	logrus.Infof("Setup of web service completed!")
	return nil
}

func (s *Server) Close() {
	s.station.Close()
}

func (s *Server) Run() {
	// Run web service
	err := s.server.Run(s.config.ListenAddress)
	if err != nil {
		logrus.Errorf("Failed to listen and serve on %s: %v", s.config.ListenAddress, err)
	}
}

func (s *Server) handleStationError(c *gin.Context, err error) bool {
	if err != nil {
		logrus.Errorf("Lost connection to DNT RoomLogg PRO: %v", err)
		s.station.Close()
		if err := s.station.Open(); err != nil {
			logrus.Errorf("Failed to restore connection to DNT RoomLogg PRO: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return true
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return true
	}
	return false
}

func (s *Server) GetCurrentData(c *gin.Context) {
	data, err := s.station.FetchCurrentData()
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetCalibrationData(c *gin.Context) {
	data, err := s.station.FetchCalibrationData()
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetIntervalMinutes(c *gin.Context) {
	data, err := s.station.FetchIntervalMinutes()
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetSettings(c *gin.Context) {
	data, err := s.station.FetchSettings()
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetAlarmSettings(c *gin.Context) {
	data, err := s.station.FetchAlarmSettings()
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetTemperatureAlarms(c *gin.Context) {
	data, err := s.station.FetchTemperatureAlarms()
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetHumidityAlarms(c *gin.Context) {
	data, err := s.station.FetchHumidityAlarms()
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) SetLanguage(c *gin.Context) {
	var input LanguageData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetLanguage(input)
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetIntervalMinutes(c *gin.Context) {
	var input IntervalData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetIntervalMinutes(input)
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetCalibration(c *gin.Context) {
	var input []*CalibrationData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetCalibrationData(input)
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetSettings(c *gin.Context) {
	var input *SettingsData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetSettings(input)
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetAlarmSettings(c *gin.Context) {
	var input *AlarmSettingsData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetAlarmSettings(input)
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetTemperatureAlarms(c *gin.Context) {
	var input []*TemperatureAlarmData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetTemperatureAlarms(input)
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetHumidityAlarms(c *gin.Context) {
	var input []*HumidityAlarmData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetHumidityAlarms(input)
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetCurrentTime(c *gin.Context) {
	err := s.station.SetTime(TimeData{time: time.Now()})
	if s.handleStationError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
