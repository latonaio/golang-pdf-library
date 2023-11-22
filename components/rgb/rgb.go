package rgb

import (
	"strconv"
	"strings"
)

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (this *RGB) ToArray() []uint8 {
	return []uint8{this.Red, this.Green, this.Blue}
}

func (this *RGB) FromArray(uints []uint8) {
	this.Red = uints[0]
	this.Green = uints[1]
	this.Blue = uints[2]
}

type Hex string

func (this *Hex) ToRGB() (*RGB, error) {
	rgb, err := Hex2RGB(this)
	return rgb, err
}

func Hex2RGB(hex *Hex) (*RGB, error) {
	rgb := RGB{}
	enhexed := strings.Replace(string(*hex), "#", "", -1)

	values, err := strconv.ParseUint(enhexed, 16, 32)

	if err != nil {
		return &rgb, err
	}

	rgb.FromArray([]uint8{uint8(values >> 16), uint8((values >> 8) & 0xFF), uint8(values & 0xFF)})

	return &rgb, nil
}
