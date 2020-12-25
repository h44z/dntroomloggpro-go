package internal

import (
	"errors"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/sirupsen/logrus"
)

type Roomlog struct {
	// Only one context should be needed for an application.  It should always be closed.
	usbOutEndpoint *gousb.OutEndpoint
	usbInEndpoint  *gousb.InEndpoint
	usbCtx         *gousb.Context
	usbDev         *gousb.Device
	usbIface       *gousb.Interface
	usbConfig      *gousb.Config
}

func NewRoomlog() *Roomlog {
	r := &Roomlog{}

	return r
}

func (r *Roomlog) Open() error {
	r.usbCtx = gousb.NewContext()

	devs, err := r.usbCtx.OpenDevices(findRoomLoggDevice())
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

	r.usbDev = devs[0]
	err = r.usbDev.SetAutoDetach(true) // release kernel driver, we need the raw input
	if err != nil {
		logrus.Errorf("Failed to set autodetach on DNT RoomLogg PRO device: %v", err)
		return err
	}

	cfgNum, err := r.usbDev.ActiveConfigNum()
	if err != nil {
		logrus.Errorf("Failed to get active default config number: %v", err)
		return err
	}
	r.usbConfig, err = r.usbDev.Config(cfgNum)
	if err != nil {
		logrus.Errorf("Failed to claim config %d of device %s: %v", cfgNum, r.usbDev, err)
		return err
	}
	r.usbIface, err = r.usbConfig.Interface(0, 0)
	if err != nil {
		logrus.Errorf("failed to select interface #%d alternate setting %d of config %d of device %s: %v", 0, 0, cfgNum, r.usbDev, err)
		return err
	}

	// Open endpoints
	inEndpointNr, outEndpointNr := getRoomLoggEndpoints(r.usbDev.Desc)
	r.usbInEndpoint, err = r.usbIface.InEndpoint(inEndpointNr)
	if err != nil {
		logrus.Errorf("%s.InEndpoint(%d): %v", r.usbIface, r.usbInEndpoint, err)
		return err
	}

	r.usbOutEndpoint, err = r.usbIface.OutEndpoint(outEndpointNr)
	if err != nil {
		logrus.Errorf("%s.OutEndpoint(%d): %v", r.usbIface, r.usbOutEndpoint, err)
		return err
	}

	return nil
}

func (r *Roomlog) Close() {
	if r.usbIface != nil {
		r.usbIface.Close()
	}
	if r.usbConfig != nil {
		r.usbConfig.Close()
	}
	if r.usbDev != nil {
		r.usbDev.Close()
	}
	if r.usbCtx != nil {
		r.usbCtx.Close()
	}
}

func (r *Roomlog) FetchData() ([]Channel, error) {
	// Request data, see: https://juergen.rocks/blog/articles/elv-raumklimastation-rs500-raspberry-pi-linux.html
	// Also a good read: https://github.com/juergen-rocks/raumklima
	outData := make([]byte, 64)
	outData[0] = 0x7b
	outData[1] = 0x03
	outData[2] = 0x40
	outData[3] = 0x7d

	logrus.Tracef("Writing: %d bytes", len(outData))
	writtenBytes, err := r.usbOutEndpoint.Write(outData)
	if err != nil || writtenBytes != len(outData) {
		logrus.Errorf("%s.Write([%d]): only %d bytes written, returned error is %v", r.usbOutEndpoint, len(outData), writtenBytes, err)
		return nil, err
	}
	logrus.Tracef("Written data to device: %v", outData)

	// Buffer large enough for 1 USB packets (64 bytes) from in-endpoint.
	logrus.Tracef("Reading: %d bytes", r.usbInEndpoint.Desc.MaxPacketSize)
	inData := make([]byte, r.usbInEndpoint.Desc.MaxPacketSize)
	// readBytes might be smaller than the buffer size. readBytes might be greater than zero even if err is not nil.
	readBytes, err := r.usbInEndpoint.Read(inData)
	if err != nil {
		logrus.Errorf("Read returned an error: %v", err)
		return nil, err
	}
	logrus.Tracef("Read %d bytes: %v", readBytes, inData)

	channels := make([]Channel, 0, 8)
	tmpChannel := Channel{}
	channelIndex := 0 // 0 indexed
	channelByte := 0

	for _, b := range inData[1:25] { // First byte is not used, each channelIndex uses 3 bytes, so 3*8 = 24 (+1 offset)
		switch channelByte {
		case 0:
			tmpChannel = Channel{
				Number: channelIndex + 1,
			}
			tmpChannel.rawTemp1 = b
			channelByte++
		case 1:
			tmpChannel.rawTemp2 = b
			channelByte++
		case 2:
			tmpChannel.rawHumidity = b
			tmpChannel.calc()
			// skip unused channels
			if tmpChannel.rawHumidity != 0xff {
				channels = append(channels, tmpChannel)
			}
			channelByte = 0
			channelIndex++
		}
	}

	return channels, nil
}

func getRoomLoggEndpoints(desc *gousb.DeviceDesc) (in int, out int) {
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

func findRoomLoggDevice() func(desc *gousb.DeviceDesc) bool {
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
