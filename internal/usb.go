package internal

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/google/gousb/usbid"

	"github.com/google/gousb"
	"github.com/sirupsen/logrus"
)

var (
	MessageStart = []byte{0x7b}
	MessageEnd   = []byte{0x40, 0x7d}
)

type UsbConnection struct {
	// Only one context should be needed for an application.  It should always be closed.
	outEndpoint *gousb.OutEndpoint
	inEndpoint  *gousb.InEndpoint
	ctx         *gousb.Context
	dev         *gousb.Device
	iface       *gousb.Interface
	config      *gousb.Config
}

func NewUsbConnection() *UsbConnection {
	c := &UsbConnection{}

	return c
}

func (c *UsbConnection) Open() error {
	c.ctx = gousb.NewContext()

	devs, err := c.ctx.OpenDevices(findUsbDevice())
	if err != nil {
		logrus.Errorf("Failed to open DNT RoomLogg PRO device: %v", err)
		return err
	}
	switch {
	case len(devs) == 0:
		logrus.Errorf("No DNT RoomLogg PRO device found.")
		return errors.New("no DNT RoomLogg PRO device found")
	case len(devs) > 1:
		logrus.Errorf("Found multiple DNT RoomLogg PRO devices, only one device is supported.")
		return errors.New("found multiple DNT RoomLogg PRO devices, only one device is supported")
	}

	c.dev = devs[0]
	err = c.dev.SetAutoDetach(true) // release kernel driver, we need the raw input
	if err != nil {
		logrus.Errorf("Failed to set autodetach on DNT RoomLogg PRO device: %v", err)
		return err
	}

	cfgNum, err := c.dev.ActiveConfigNum()
	if err != nil {
		logrus.Errorf("Failed to get active default config number: %v", err)
		return err
	}
	c.config, err = c.dev.Config(cfgNum)
	if err != nil {
		logrus.Errorf("Failed to claim config %d of device %s: %v", cfgNum, c.dev, err)
		return err
	}
	c.iface, err = c.config.Interface(0, 0)
	if err != nil {
		logrus.Errorf("Failed to select interface #%d alternate setting %d of config %d of device %s: %v", 0, 0, cfgNum, c.dev, err)
		return err
	}

	// Open endpoints
	inEndpointNr, outEndpointNr := getUsbEndpoints(c.dev.Desc)
	c.inEndpoint, err = c.iface.InEndpoint(inEndpointNr)
	if err != nil {
		logrus.Errorf("Failed to get in-endpoint: %v", err)
		return err
	}

	c.outEndpoint, err = c.iface.OutEndpoint(outEndpointNr)
	if err != nil {
		logrus.Errorf("Failed to get out-endpoint: %v", err)
		return err
	}

	return nil
}

func (c *UsbConnection) Close() {
	if c.iface != nil {
		c.iface.Close()
	}
	if c.config != nil {
		c.config.Close()
	}
	if c.dev != nil {
		c.dev.Close()
	}
	if c.ctx != nil {
		c.ctx.Close()
	}
}

func (c *UsbConnection) rawRequest(body []byte, noResponse ...bool) ([]byte, error) {
	lenBodyBytes := len(body)
	lenPaddedBytes := c.outEndpoint.Desc.MaxPacketSize
	// Apply padding
	if lenBodyBytes < lenPaddedBytes {
		lenPadding := lenPaddedBytes - lenBodyBytes
		for i := 0; i < lenPadding; i++ {
			body = append(body, 0)
		}
	}

	logrus.Tracef("Writing raw message: %d bytes, %d bytes with padding", lenBodyBytes, lenPaddedBytes)
	numWrittenBytes, err := c.outEndpoint.Write(body)
	if err != nil || numWrittenBytes != lenPaddedBytes {
		logrus.Errorf("Failed to write raw bytes: %v. Written: %d bytes", err, numWrittenBytes)
		return nil, fmt.Errorf("failed to write raw bytes: %v, written: %d bytes", err, numWrittenBytes)
	}
	logrus.Tracef("Written data to device: %d bytes", numWrittenBytes)

	if noResponse != nil && len(noResponse) > 0 && noResponse[0] == true {
		return nil, nil
	}

	// Buffer large enough for 1 USB packets (64 bytes) from in-endpoint.
	logrus.Tracef("Reading raw message: %d bytes", c.inEndpoint.Desc.MaxPacketSize)
	inData := make([]byte, c.inEndpoint.Desc.MaxPacketSize)
	// numReadBytes might be smaller than the buffer size. numReadBytes might be greater than zero even if err is not nil.
	numReadBytes, err := c.inEndpoint.Read(inData)
	if err != nil {
		logrus.Errorf("Failed to read raw bytes: %v", err)
		return nil, fmt.Errorf("failed to read raw bytes: %v, read %d bytes", err, numReadBytes)
	}
	logrus.Tracef("Read %d bytes from device", numReadBytes)

	return inData, nil
}

func (c *UsbConnection) Request(command byte, payload []byte, noResponse ...bool) ([]byte, error) {
	var body = make([]byte, 0, 4+len(payload)) // 4 for command and start/end bytes
	body = append(body, MessageStart...)
	body = append(body, command)
	body = append(body, payload...)
	body = append(body, MessageEnd...)

	return c.rawRequest(body, noResponse...)
}

func getUsbEndpoints(desc *gousb.DeviceDesc) (in int, out int) {
	for _, cfg := range desc.Configs { // There should only be one config
		for _, intf := range cfg.Interfaces { // There should only be one interface
			for _, ifSetting := range intf.AltSettings { // There should only be one alt-setting
				for _, end := range ifSetting.Endpoints {
					switch end.Direction {
					case gousb.EndpointDirectionOut:
						out = end.Number
					case gousb.EndpointDirectionIn:
						in = end.Number
					}
				}
			}
		}
	}
	return
}

func findUsbDevice() func(desc *gousb.DeviceDesc) bool {
	return func(desc *gousb.DeviceDesc) bool {

		// The usbid package can be used to print out human readable information.
		logrus.Debugf("Checking USB device: %03d.%03d %s:%s %s", desc.Bus, desc.Address, desc.Vendor, desc.Product, usbid.Describe(desc))

		if desc.Product == gousb.ID(0x5750) && desc.Vendor == gousb.ID(0x0483) {
			logrus.Infof("Found DNT RoomLogg PRO on USB bus %03d.%03d", desc.Bus, desc.Address)

			return true
		}

		return false
	}
}

func GetMessagePayload(raw []byte) ([]byte, error) {
	if raw == nil || len(raw) < 3 {
		return nil, errors.New("invalid raw message size")
	}

	startIndex := 1 // First byte can be dismissed, its 0x7b (MessageStart)
	endIndex := bytes.Index(raw, MessageEnd)

	if endIndex == -1 {
		return nil, errors.New("unable to find message end")
	}

	return raw[startIndex:endIndex], nil
}
