package internal

import (
	"encoding/binary"
)

type Channel struct {
	Number      int
	Temperature float64
	Humidity    float64

	rawTemp1    byte
	rawTemp2    byte
	rawHumidity byte
}

func (c *Channel) calc() {
	tmp := int16(binary.BigEndian.Uint16([]byte{c.rawTemp1, c.rawTemp2}))
	c.Temperature = float64(tmp) / 10
	c.Humidity = float64(c.rawHumidity)
}
