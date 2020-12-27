package internal

import (
	"encoding/binary"
)

type TimeFormat uint8
type DateFormat uint8
type GraphType uint8
type GraphInterval uint8
type Unit uint8
type Flag uint8
type Language uint8

const (
	LanguageDE Language = 0x00
	LanguageEN Language = 0x01

	DSTOn    Flag = 0x01
	DSTOff   Flag = 0x00
	AlarmOn  Flag = 0x00
	AlarmOff Flag = 0x01

	TimeFormatEurope        TimeFormat = 0x00
	TimeFormatEnglishPrefix TimeFormat = 0x01
	TimeFormatEnglishSuffix TimeFormat = 0x02

	DateFormatYYYYMMDD DateFormat = 0x00
	DateFormatMMDDYYYY DateFormat = 0x01
	DateFormatDDMMYYYY DateFormat = 0x02

	GraphInterval12h GraphInterval = 0x0c
	GraphInterval24h GraphInterval = 0x18
	GraphInterval48h GraphInterval = 0x30
	GraphInterval72h GraphInterval = 0x48

	GraphTypeTemperature GraphType = 0x00
	GraphTypeHumidity    GraphType = 0x01
	GraphTypeDewPoint    GraphType = 0x02
	GraphTypeHeatIndex   GraphType = 0x03

	UnitCelsius    Unit = 0x00
	UnitFahrenheit Unit = 0x01

	CommandGetAlarmSettings    = 0x06
	CommandGetTemperatureAlarm = 0x08
	CommandGetHumidityAlarm    = 0x09
	CommandGetCalibration      = 0x05
	CommandGetCurrentData      = 0x03
	CommandGetSettings         = 0x04
	CommandGetInterval         = 0x41

	CommandSetAlarmSettings    = 0x12
	CommandSetTemperatureAlarm = 0x14
	CommandSetHumidityAlarm    = 0x15
	CommandSetCalibration      = 0x11
	CommandSetSettings         = 0x10
	CommandSetInterval         = 0x40
	CommandSetLanguage         = 0x0b
	CommandSetTime             = 0x30

	CommandStartStore = 0x51
	CommandEndStore   = 0x2f
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

type IntervalData uint8 // Minutes

func NewIntervalData(raw []byte) IntervalData {
	return IntervalData(raw[0])
}

func (d IntervalData) RawBytes() []byte {
	r := make([]byte, 1)

	r[0] = byte(d)

	return r
}

type LanguageData Language // Either DE or EN

func NewLanguageData(raw []byte) LanguageData {
	return LanguageData(raw[0])
}

func (d LanguageData) RawBytes() []byte {
	r := make([]byte, 1)

	r[0] = byte(d)

	return r
}

type SettingsAreaData struct {
	Area int

	// Temperature and Humidity alarm can be disabled/enabled for each channel and also for high/low values.
	// Map index is channel number (0 based)
	Temperature map[uint8]bool
	DewPoint    map[uint8]bool
	HeatIndex   map[uint8]bool
}

func NewSettingsAreaData(raw []byte) *SettingsAreaData {
	d := &SettingsAreaData{}
	d.Temperature = make(map[uint8]bool, 5)
	d.DewPoint = make(map[uint8]bool, 5)
	d.HeatIndex = make(map[uint8]bool, 5)

	// First bit = bit on the right: First bit = temperature, second bit = dewpoint, third bit = heatindex, fourth bit = temp of ch2, ....
	areaData := binary.BigEndian.Uint32(raw)

	// Test bits
	channel := uint8(0)
	bitType := 0
	for i := uint(0); i < 3*8; i++ { // 3 bits for each channel, 8 channels max
		switch bitType {
		case 0:
			d.Temperature[channel] = hasBitSet32(areaData, i)
		case 1:
			d.DewPoint[channel] = hasBitSet32(areaData, i)
		case 2:
			d.HeatIndex[channel] = hasBitSet32(areaData, i)
		}

		bitType++
		if bitType%3 == 0 {
			bitType = 0
			channel++
		}
	}

	return d
}

func (d *SettingsAreaData) RawBytes() []byte {
	r := make([]byte, 4)

	areaData := uint32(0)
	channel := uint8(0)
	bitType := 0
	for i := uint(0); i < 3*8; i++ { // 3 bits for each channel, 8 channels max
		switch bitType {
		case 0:
			if d.Temperature[channel] {
				areaData = setBit32(areaData, i)
			}
		case 1:
			if d.DewPoint[channel] {
				areaData = setBit32(areaData, i)
			}
		case 2:
			if d.HeatIndex[channel] {
				areaData = setBit32(areaData, i)
			}
		}

		bitType++
		if bitType%3 == 0 {
			bitType = 0
			channel++
		}
	}

	binary.BigEndian.PutUint32(r, areaData)

	return r
}

func setBit32(n uint32, pos uint) uint32 {
	n |= 1 << pos
	return n
}

func hasBitSet32(n uint32, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}

type SettingsData struct {
	GraphType     GraphType
	GraphInterval GraphInterval
	TimeFormat    TimeFormat
	DateFormat    DateFormat
	DST           Flag
	TimeZone      int8
	Units         Unit
	Areas         [5]*SettingsAreaData
}

func NewSettingsData(raw []byte) *SettingsData {
	d := &SettingsData{}
	d.GraphType = GraphType(raw[0])
	d.GraphInterval = GraphInterval(raw[1])
	d.TimeFormat = TimeFormat(raw[2])
	d.DateFormat = DateFormat(raw[3])
	d.DST = Flag(raw[4])
	d.TimeZone = int8(raw[5])
	d.Units = Unit(raw[6])

	areaIndex := 0 // index is 0 based, area number is index+1
	iOffset := 7
	for i := 0; i < 4*5; i += 4 { // One area consists of 4 bytes, 5 areas in total
		area := NewSettingsAreaData(raw[i+iOffset : i+iOffset+4])
		area.Area = areaIndex + 1
		d.Areas[areaIndex] = area
		areaIndex++
	}

	return d
}

func (d *SettingsData) RawBytes() []byte {
	r := make([]byte, 27)
	r[0] = byte(d.GraphType)
	r[1] = byte(d.GraphInterval)
	r[2] = byte(d.TimeFormat)
	r[3] = byte(d.DateFormat)
	r[4] = byte(d.DST)
	r[5] = byte(d.TimeZone)
	r[6] = byte(d.Units)

	areaIndex := 7
	for i := 0; i < 5; i++ {
		areaRawBytes := d.Areas[i].RawBytes()
		for _, b := range areaRawBytes {
			r[areaIndex] = b
			areaIndex++
		}
	}

	return r
}

type AlarmSettingsData struct {
	EnableTemperatureAlarm Flag
	EnableHumidityAlarm    Flag

	// Map index = channel number (0 based)
	TemperatureLowAlarm  map[uint8]bool
	TemperatureHighAlarm map[uint8]bool
	HumidityLowAlarm     map[uint8]bool
	HumidityHighAlarm    map[uint8]bool
}

func NewAlarmSettingsData(raw []byte) *AlarmSettingsData {
	d := &AlarmSettingsData{}
	d.EnableTemperatureAlarm = Flag(raw[0]) // 0x01 = off, 0x00 = on
	d.EnableHumidityAlarm = Flag(raw[1])

	// check remaining 4 bytes, first byte = hum high, second byte = hum low, third byte = tmp high, fourth byte = tmp low
	// ch1 = bit 0, ch2 = bit 1, ch3 = bit 2, ... lowest bit is on the right
	for i := 2; i < 6; i++ {
		flags := make(map[uint8]bool, 8)
		for b := uint8(0); b < 8; b++ {
			flags[b] = hasBitSet8(raw[i], uint(b))
		}

		switch i {
		case 2: // hum high flags
			d.HumidityHighAlarm = flags
		case 3: // hum low flags
			d.HumidityLowAlarm = flags
		case 4: // tmp high flags
			d.TemperatureHighAlarm = flags
		case 5: // tmp low flags
			d.TemperatureLowAlarm = flags
		}
	}
	return d
}

func setBit8(n uint8, pos uint) uint8 {
	n |= 1 << pos
	return n
}

func hasBitSet8(n uint8, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}

func (d *AlarmSettingsData) RawBytes() []byte {
	r := make([]byte, 6)
	r[0] = byte(d.EnableTemperatureAlarm)
	r[1] = byte(d.EnableHumidityAlarm)
	r[2] = byte(0)
	r[3] = byte(0)
	r[4] = byte(0)
	r[5] = byte(0)

	for i := uint8(0); i < 8; i++ {
		if d.HumidityHighAlarm[i] {
			r[2] = setBit8(r[2], uint(i))
		}
		if d.HumidityLowAlarm[i] {
			r[3] = setBit8(r[3], uint(i))
		}
		if d.TemperatureHighAlarm[i] {
			r[4] = setBit8(r[4], uint(i))
		}
		if d.TemperatureLowAlarm[i] {
			r[5] = setBit8(r[5], uint(i))
		}
	}

	return r
}

type HumidityAlarmData struct {
	Channel int
	Low     float64
	High    float64
}

func NewHumidityAlarmData(raw []byte) *HumidityAlarmData {
	d := &HumidityAlarmData{}
	d.High = float64(raw[0])
	d.Low = float64(raw[1])

	return d
}

func NewHumidityAlarmsData(raw []byte) []*HumidityAlarmData {
	numChannels := len(raw) / 2 // 1 channel has 2 bytes

	channels := make([]*HumidityAlarmData, 0, numChannels)
	channelIndex := 1
	for i := 0; i < 2*numChannels; i += 2 {
		channel := NewHumidityAlarmData(raw[i : i+2])
		channel.Channel = channelIndex
		channels = append(channels, channel)
		channelIndex++
	}

	return channels
}

func (d *HumidityAlarmData) RawBytes() []byte {
	r := make([]byte, 2)
	r[0] = byte(d.High)
	r[1] = byte(d.Low)

	return r
}

type TemperatureAlarmData struct {
	Channel int
	Low     float64
	High    float64
}

func NewTemperatureAlarmData(raw []byte) *TemperatureAlarmData {
	d := &TemperatureAlarmData{}
	tmp := int16(binary.BigEndian.Uint16([]byte{raw[0], raw[1]}))
	d.High = float64(tmp) / 10
	tmp = int16(binary.BigEndian.Uint16([]byte{raw[2], raw[3]}))
	d.Low = float64(tmp) / 10

	return d
}

func NewTemperatureAlarmsData(raw []byte) []*TemperatureAlarmData {
	numChannels := len(raw) / 4 // 1 channel has 4 bytes

	channels := make([]*TemperatureAlarmData, 0, numChannels)
	channelIndex := 1
	for i := 0; i < 4*numChannels; i += 4 {
		channel := NewTemperatureAlarmData(raw[i : i+4])
		channel.Channel = channelIndex
		channels = append(channels, channel)
		channelIndex++
	}

	return channels
}

func (d *TemperatureAlarmData) RawBytes() []byte {
	r := make([]byte, 4)
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(d.High*10))
	r[0] = tmp[0]
	r[1] = tmp[1]
	binary.BigEndian.PutUint16(tmp, uint16(d.Low*10))
	r[2] = tmp[0]
	r[3] = tmp[1]

	return r
}

type CalibrationData struct {
	Channel     int
	Temperature float64
	Humidity    float64
}

func NewCalibrationData(raw []byte) *CalibrationData {
	d := &CalibrationData{}
	tmp := int16(binary.BigEndian.Uint16([]byte{raw[0], raw[1]}))
	d.Temperature = float64(tmp) / 10
	d.Humidity = float64(raw[2])

	return d
}

func NewCalibrationsData(raw []byte) []*CalibrationData {
	numChannels := len(raw) / 3 // 1 channel has 3 bytes
	channels := make([]*CalibrationData, 0, numChannels)
	channelIndex := 1
	for i := 0; i < 3*numChannels; i += 3 {
		channel := NewCalibrationData(raw[i : i+3])
		channel.Channel = channelIndex
		channels = append(channels, channel)
		channelIndex++
	}

	return channels
}

func (d *CalibrationData) RawBytes() []byte {
	r := make([]byte, 3)

	r[2] = byte(d.Humidity)
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(d.Temperature*10))
	r[0] = tmp[0]
	r[1] = tmp[1]

	return r
}
