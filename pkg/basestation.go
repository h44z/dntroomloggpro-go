package pkg

import (
	"errors"
	"fmt"

	"github.com/h44z/dntroomloggpro-go/internal"
	"github.com/sirupsen/logrus"
)

type RoomLogg struct {
	// Only one context should be needed for an application.  It should always be closed.
	cfg *RoomLoggConfig
	usb *internal.UsbConnection
}

func NewRoomLogg(cfg *RoomLoggConfig) *RoomLogg {
	r := &RoomLogg{cfg: cfg}
	r.usb = internal.NewUsbConnection()

	return r
}

func (r *RoomLogg) Open() error {
	return r.usb.Open()
}

func (r *RoomLogg) Close() {
	r.usb.Close()
}

func (r *RoomLogg) Reconnect() error {
	r.Close()
	err := r.Open()
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomLogg) FetchCurrentData() ([]*ChannelData, error) { // Returns already calibrated data
	dataBytes, err := r.usb.Request(CommandGetCurrentData, nil)
	if err != nil || dataBytes[0] != internal.MessageStart[0] {
		logrus.Errorf("Failed to fetch current data: %v", err)
		return nil, err
	}

	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewChannelsData(payload), nil
}

func (r *RoomLogg) FetchCalibrationData() ([]*CalibrationData, error) {
	dataBytes, err := r.usb.Request(CommandGetCalibration, nil)
	if err != nil || dataBytes[0] != internal.MessageStart[0] {
		logrus.Errorf("Failed to fetch calibration data: %v", err)
		return nil, err
	}

	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewCalibrationsData(payload), nil
}

func (r *RoomLogg) FetchIntervalMinutes() (IntervalData, error) {
	dataBytes, err := r.usb.Request(CommandGetInterval, nil)
	if err != nil || dataBytes[0] != internal.MessageStart[0] {
		logrus.Errorf("Failed to fetch interval data: %v", err)
		return 0, err
	}

	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return 0, err
	}

	return NewIntervalData(payload), nil
}

func (r *RoomLogg) FetchSettings() (*SettingsData, error) {
	dataBytes, err := r.usb.Request(CommandGetSettings, nil)
	if err != nil || dataBytes[0] != internal.MessageStart[0] {
		logrus.Errorf("Failed to fetch settings data: %v", err)
		return nil, err
	}

	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewSettingsData(payload), nil
}

func (r *RoomLogg) FetchAlarmSettings() (*AlarmSettingsData, error) {
	dataBytes, err := r.usb.Request(CommandGetAlarmSettings, nil)
	if err != nil || dataBytes[0] != internal.MessageStart[0] {
		logrus.Errorf("Failed to fetch alarm settings data: %v", err)
		return nil, err
	}

	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewAlarmSettingsData(payload), nil
}

func (r *RoomLogg) FetchTemperatureAlarms() ([]*TemperatureAlarmData, error) {
	dataBytes, err := r.usb.Request(CommandGetTemperatureAlarm, nil)
	if err != nil || dataBytes[0] != internal.MessageStart[0] {
		logrus.Errorf("Failed to fetch temperature alarm data: %v", err)
		return nil, err
	}

	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewTemperatureAlarmsData(payload), nil
}

func (r *RoomLogg) FetchHumidityAlarms() ([]*HumidityAlarmData, error) {
	dataBytes, err := r.usb.Request(CommandGetHumidityAlarm, nil)
	if err != nil || dataBytes[0] != internal.MessageStart[0] {
		logrus.Errorf("Failed to fetch humidity alarm data: %v", err)
		return nil, err
	}

	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewHumidityAlarmsData(payload), nil
}

func (r *RoomLogg) SetIntervalMinutes(minutes IntervalData) error {
	if minutes < 0 || minutes > 240 {
		return errors.New("value out of range")
	}
	// Interval sync has no start-store command

	if _, err := r.usb.Request(CommandSetInterval, minutes.RawBytes(), true); err != nil {
		logrus.Errorf("Failed to set interval data: %v", err)
		return err
	}

	if err := r.endStore(); err != nil {
		return err
	}

	return nil
}

func (r *RoomLogg) SetLanguage(lang LanguageData) error {
	if lang < 0 || lang > 1 {
		return errors.New("value out of range")
	}

	if err := r.startStore(); err != nil {
		return err
	}

	if _, err := r.usb.Request(CommandSetLanguage, lang.RawBytes(), true); err != nil {
		logrus.Errorf("Failed to set language data: %v", err)
		return err
	}

	if err := r.endStore(); err != nil {
		return err
	}

	return nil
}

func (r *RoomLogg) SetTime(time TimeData) error {
	if err := r.startStore(); err != nil {
		return err
	}

	if _, err := r.usb.Request(CommandSetTime, time.RawBytes(), true); err != nil {
		logrus.Errorf("Failed to set time data: %v", err)
		return err
	}

	// Time sync has no end-store command

	return nil
}

func (r *RoomLogg) SetCalibrationData(calibration []*CalibrationData) error {
	if err := r.startStore(); err != nil {
		return err
	}

	rawBytes := make([]byte, 0, 24) // 3 * 8 bytes
	for i := 0; i < 8; i++ {
		rawBytes = append(rawBytes, calibration[i].RawBytes()...)
	}
	if _, err := r.usb.Request(CommandSetCalibration, rawBytes, true); err != nil {
		logrus.Errorf("Failed to set calibration data: %v", err)
		return err
	}

	if err := r.endStore(); err != nil {
		return err
	}

	return nil
}

func (r *RoomLogg) SetSettings(settings *SettingsData) error {
	if err := r.startStore(); err != nil {
		return err
	}

	if _, err := r.usb.Request(CommandSetSettings, settings.RawBytes(), true); err != nil {
		logrus.Errorf("Failed to set settings data: %v", err)
		return err
	}

	if err := r.endStore(); err != nil {
		return err
	}

	return nil
}

func (r *RoomLogg) SetAlarmSettings(settings *AlarmSettingsData) error {
	if err := r.startStore(); err != nil {
		return err
	}

	if _, err := r.usb.Request(CommandSetAlarmSettings, settings.RawBytes(), true); err != nil {
		logrus.Errorf("Failed to set alarm settings data: %v", err)
		return err
	}

	if err := r.endStore(); err != nil {
		return err
	}

	return nil
}

func (r *RoomLogg) SetTemperatureAlarms(alarms []*TemperatureAlarmData) error {
	if err := r.startStore(); err != nil {
		return err
	}

	rawBytes := make([]byte, 0, 32) // 4 * 8 bytes
	for i := 0; i < 8; i++ {
		rawBytes = append(rawBytes, alarms[i].RawBytes()...)
	}
	if _, err := r.usb.Request(CommandSetTemperatureAlarm, rawBytes, true); err != nil {
		logrus.Errorf("Failed to set temperature alarm data: %v", err)
		return err
	}

	if err := r.endStore(); err != nil {
		return err
	}

	return nil
}

func (r *RoomLogg) SetHumidityAlarms(alarms []*HumidityAlarmData) error {
	if err := r.startStore(); err != nil {
		return err
	}

	rawBytes := make([]byte, 0, 8)
	for i := 0; i < 8; i++ {
		rawBytes = append(rawBytes, alarms[i].RawBytes()...)
	}
	if _, err := r.usb.Request(CommandSetHumidityAlarm, rawBytes, true); err != nil {
		logrus.Errorf("Failed to set humidity alarm data: %v", err)
		return err
	}

	if err := r.endStore(); err != nil {
		return err
	}

	return nil
}

func (r *RoomLogg) startStore() error {
	dataBytes, err := r.usb.Request(CommandStartStore, nil)
	if err != nil {
		logrus.Errorf("Failed to start store command: %v", err)
		return err
	}
	payload, err := internal.GetMessagePayload(dataBytes)
	if err != nil || len(payload) != 1 || payload[0] != 0x00 {
		logrus.Errorf("Failed to start store command, bad response: %v, %v", err, payload)
		return fmt.Errorf("failed to start store command, bad response: %v, %v", err, payload)
	}

	return nil
}

func (r *RoomLogg) endStore() error {
	if _, err := r.usb.Request(CommandEndStore, nil, true); err != nil {
		logrus.Errorf("Failed to end store command: %v", err)
		return err
	}

	return nil
}
