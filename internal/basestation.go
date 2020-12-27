package internal

import (
	"github.com/sirupsen/logrus"
)

type RoomLogg struct {
	// Only one context should be needed for an application.  It should always be closed.
	usb *UsbConnection
}

func NewRoomLogg() *RoomLogg {
	r := &RoomLogg{}
	r.usb = NewUsbConnection()

	return r
}

func (r *RoomLogg) Open() error {
	return r.usb.Open()
}

func (r *RoomLogg) Close() {
	r.usb.Close()
}

func (r *RoomLogg) FetchCurrentData() ([]*ChannelData, error) { // Returns already calibrated data
	dataBytes, err := r.usb.Request(CommandGetCurrentData, nil)
	if err != nil || dataBytes[0] != MessageStart[0] {
		logrus.Errorf("Failed to fetch current data: %v", err)
		return nil, err
	}

	payload, err := GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewChannelsData(payload), nil
}

func (r *RoomLogg) FetchCalibrationData() ([]*CalibrationData, error) {
	dataBytes, err := r.usb.Request(CommandGetCalibration, nil)
	if err != nil || dataBytes[0] != MessageStart[0] {
		logrus.Errorf("Failed to fetch calibration data: %v", err)
		return nil, err
	}

	payload, err := GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewCalibrationsData(payload), nil
}

func (r *RoomLogg) FetchIntervalMinutes() (IntervalData, error) {
	dataBytes, err := r.usb.Request(CommandGetInterval, nil)
	if err != nil || dataBytes[0] != MessageStart[0] {
		logrus.Errorf("Failed to fetch interval data: %v", err)
		return 0, err
	}

	payload, err := GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return 0, err
	}

	return NewIntervalData(payload), nil
}

func (r *RoomLogg) FetchSettings() (*SettingsData, error) {
	dataBytes, err := r.usb.Request(CommandGetSettings, nil)
	if err != nil || dataBytes[0] != MessageStart[0] {
		logrus.Errorf("Failed to fetch settings data: %v", err)
		return nil, err
	}

	payload, err := GetMessagePayload(dataBytes)
	if err != nil {
		logrus.Errorf("Failed to validate message payload: %v", err)
		return nil, err
	}

	return NewSettingsData(payload), nil
}
