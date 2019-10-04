package ad9850

import (
	"github.com/stianeikeland/go-rpio"
)

type DDS struct {
	Frequency int
	Enabled bool
	pinCLK rpio.Pin
	pinFQUD rpio.Pin
	pinDATA rpio.Pin
	pinRST rpio.Pin
}

func Init(clk uint8, fqud uint8, data uint8, rst uint8) (err error, dds DDS) {

	err = rpio.Open()
	if err != nil {
		return err, dds
	}

	dds.pinCLK = rpio.Pin(clk)
	dds.pinFQUD = rpio.Pin(fqud)
	dds.pinDATA = rpio.Pin(data)
	dds.pinRST = rpio.Pin(rst)

	dds.pinCLK.Output()
	dds.pinFQUD.Output()
	dds.pinDATA.Output()
	dds.pinRST.Output()

	dds.Reset()
	dds.Disable()

	return err, dds

}

func (dds *DDS) Reset() {

	dds.pinCLK.Low()
	dds.pinFQUD.Low()
	dds.pinDATA.Low()
	dds.pinRST.Low()

	dds.pinRST.High()
	dds.pinRST.Low()

	dds.pinCLK.High()
	dds.pinCLK.Low()

	dds.pinFQUD.High()
	dds.pinFQUD.Low()

}

func (dds *DDS) sendByte(b uint) {
	var i uint
	for i = 0; i < 8; i++ {
		if b&(1<<i) > 0 {
			dds.pinDATA.High()
		} else {
			dds.pinDATA.Low()
		}
		dds.pinCLK.High()
		dds.pinCLK.Low()
	}
}

func (dds *DDS) sendBytes() {
	var b uint
	n1 := float64(dds.Frequency)
	n2 := float64(4294967295)
	n3 := float64(125000000)
	f := uint(n1 * (n2 / n3))
	for b = 0; b < 4; b++ {
		dds.sendByte((f>>(8*b))&0xFF)
	}
	if (dds.Enabled) {
		dds.sendByte(0)
	} else {
		dds.sendByte(1<<2)
	}
	dds.pinFQUD.High()
	dds.pinFQUD.Low()
}

func (dds *DDS) SetFrequency(frequency int) {
	dds.Frequency = frequency
	dds.sendBytes()
}

func (dds *DDS) Disable() {
	dds.Enabled = false
	dds.sendBytes()
}

func (dds *DDS) Enable() {
	dds.Enabled = true
	dds.sendBytes()
}