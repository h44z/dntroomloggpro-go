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

func (r *RoomLogg) FetchCurrentData() ([]*ChannelData, error) {
	dataBytes, err := r.usb.Request(0x03, nil)
	if err != nil || dataBytes[0] != MessageStart[0] {
		logrus.Errorf("Failed to fetch current data: %v", err)
		return nil, err
	}

	return NewChannelsData(dataBytes[1:25]), nil // First byte is not used, each channelIndex uses 3 bytes, so 3*8 = 24 (+1 offset)
}
