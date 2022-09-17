package pkg

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Server struct {
	// Core components
	cfg    *RestConfig
	server *gin.Engine

	// cache
	settings *SettingsData
	channels []*ChannelData
	isOnline bool
	mux      sync.RWMutex

	// direct fetching
	station *RoomLogg
}

func getExecutableDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Errorf("[REST] Failed to get executable directory: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "assets")); os.IsNotExist(err) {
		return "." // assets directory not found -> we are developing in goland =)
	}

	return dir
}

func NewServer(cfg *RestConfig) (*Server, error) {
	s := &Server{cfg: cfg}

	err := s.Setup()

	return s, err
}

func (s *Server) Setup() error {
	dir := getExecutableDirectory()
	rDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	logrus.Infof("[REST] Real working directory: %s", rDir)
	logrus.Infof("[REST] Current working directory: %s", dir)

	// Init rand
	rand.Seed(time.Now().UnixNano())

	// Setup http server
	s.server = gin.Default()

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

	logrus.Infof("[REST] Setup of web service completed!")
	return nil
}

func (s *Server) Publish(settings *SettingsData, channels []*ChannelData, isOnline bool) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.settings = settings
	s.channels = channels
	s.isOnline = isOnline

	return nil
}

func (s *Server) Run() {
	// Run web service
	err := s.server.Run(s.cfg.ListenAddress)
	if err != nil {
		logrus.Errorf("[REST] Failed to listen and serve on %s: %v", s.cfg.ListenAddress, err)
	}
}

func (s *Server) GetCurrentData(c *gin.Context) {
	s.mux.RLock()
	defer s.mux.RLock()

	if !s.isOnline {
		c.Status(http.StatusGatewayTimeout)
		return
	}

	c.JSON(http.StatusOK, s.channels)
}

func (s *Server) GetSettings(c *gin.Context) {
	s.mux.RLock()
	defer s.mux.RLock()

	if !s.isOnline {
		c.Status(http.StatusGatewayTimeout)
		return
	}

	c.JSON(http.StatusOK, s.settings)
}

func (s *Server) SetRoomLogInstance(station *RoomLogg) {
	s.station = station
}

// special functions

func (s *Server) GetCalibrationData(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	data, err := s.station.FetchCalibrationData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetIntervalMinutes(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	data, err := s.station.FetchIntervalMinutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetAlarmSettings(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	data, err := s.station.FetchAlarmSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetTemperatureAlarms(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	data, err := s.station.FetchTemperatureAlarms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetHumidityAlarms(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	data, err := s.station.FetchHumidityAlarms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) SetLanguage(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	var input LanguageData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetLanguage(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetIntervalMinutes(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	var input IntervalData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetIntervalMinutes(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetCalibration(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	var input []*CalibrationData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetCalibrationData(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetSettings(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	var input *SettingsData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetSettings(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetAlarmSettings(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	var input *AlarmSettingsData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetAlarmSettings(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetTemperatureAlarms(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	var input []*TemperatureAlarmData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetTemperatureAlarms(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetHumidityAlarms(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	var input []*HumidityAlarmData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.station.SetHumidityAlarms(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (s *Server) SetCurrentTime(c *gin.Context) {
	if s.station == nil {
		c.Status(http.StatusMethodNotAllowed)
		return
	}

	err := s.station.SetTime(TimeData{time: time.Now()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
