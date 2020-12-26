package internal

import (
	"encoding/binary"
)

const (
	LanguageDE              = 0x01
	LanguageEN              = 0x00
	DSTOn                   = 0x01
	DSTOff                  = 0x00
	TimeFormatEurope        = 0x00
	TimeFormatEnglishPrefix = 0x01
	TimeFormatEnglishSuffix = 0x02
	DateFormatYYYYMMDD      = 0x00
	DateFormatMMDDYYYY      = 0x01
	DateFormatDDMMYYYY      = 0x02
	GraphInterval12h        = 0x0c
	GraphInterval24h        = 0x18
	GraphInterval48h        = 0x30
	GraphInterval72h        = 0x48
	GraphTypeTemperature    = 0x01
	GraphTypeHumidity       = 0x02
	GraphTypeDewpoint       = 0x03
	GraphTypeHeatindex      = 0x04
	UnitCelsius             = 0x00
	UnitFahrenheit          = 0x01
)

type Data interface {
	RawBytes() []byte
}

type ChannelData struct {
	Number      int
	Temperature float64
	Humidity    float64
}

func NewChannelData(raw []byte) *ChannelData {
	d := &ChannelData{}
	tmp := int16(binary.BigEndian.Uint16([]byte{raw[0], raw[1]}))
	d.Temperature = float64(tmp) / 10
	d.Humidity = float64(raw[2])

	return d
}

func NewChannelsData(raw []byte) []*ChannelData {
	numChannels := len(raw) / 3 // 1 channel has 3 bytes
	channels := make([]*ChannelData, 0, numChannels)
	channelIndex := 1
	for i := 0; i < 3*numChannels; i += 3 {
		channel := NewChannelData(raw[i : i+3])
		channel.Number = channelIndex
		if channel.Humidity != 255 {
			channels = append(channels, channel)
		}
		channelIndex++
	}

	return channels
}

func (d *ChannelData) RawBytes() []byte {
	r := make([]byte, 3)

	r[2] = byte(d.Humidity)
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(d.Temperature*10))
	r[0] = tmp[0]
	r[1] = tmp[1]

	return r
}

type IntervalData int // Minutes

func NewIntervalData(raw []byte) IntervalData {
	return IntervalData(raw[0])
}

func (d IntervalData) RawBytes() []byte {
	r := make([]byte, 1)

	r[0] = byte(d)

	return r
}

type LanguageData int // Either DE or EN

func NewLanguageData(raw []byte) LanguageData {
	return LanguageData(raw[0])
}

func (d LanguageData) RawBytes() []byte {
	r := make([]byte, 1)

	r[0] = byte(d)

	return r
}

type SettingsData struct {
}

func NewSettingsData(raw []byte) *SettingsData {
	return &SettingsData{}
}

func (d *SettingsData) RawBytes() []byte {
	r := make([]byte, 1)

	return r
}

type AlarmSettingsData struct {
}

func NewAlarmSettingsData(raw []byte) *AlarmSettingsData {
	return &AlarmSettingsData{}
}

func (d *AlarmSettingsData) RawBytes() []byte {
	r := make([]byte, 1)

	return r
}

type HumidityAlarmData struct {
}

func NewHumidityAlarmData(raw []byte) *HumidityAlarmData {
	return &HumidityAlarmData{}
}

func NewHumidityAlarmsData(raw []byte) []*HumidityAlarmData {
	numChannels := len(raw) / 2 // 1 channel has 2 bytes
	return make([]*HumidityAlarmData, 0, numChannels)
}

func (d *HumidityAlarmData) RawBytes() []byte {
	r := make([]byte, 1)

	return r
}

type TemperatureAlarmData struct {
}

func NewTemperatureAlarmData(raw []byte) *TemperatureAlarmData {
	return &TemperatureAlarmData{}
}

func NewTemperatureAlarmsData(raw []byte) []*TemperatureAlarmData {
	numChannels := len(raw) / 4 // 1 channel has 4 bytes
	return make([]*TemperatureAlarmData, 0, numChannels)
}

func (d *TemperatureAlarmData) RawBytes() []byte {
	r := make([]byte, 1)

	return r
}

type CalibrationData struct {
}

func NewCalibrationData(raw []byte) *CalibrationData {
	return &CalibrationData{}
}

func CalibrationsData(raw []byte) []*CalibrationData {
	numChannels := len(raw) / 3 // 1 channel has 3 bytes
	return make([]*CalibrationData, 0, numChannels)
}

func (d *CalibrationData) RawBytes() []byte {
	r := make([]byte, 1)

	return r
}
