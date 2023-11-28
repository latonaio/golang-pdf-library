package lnpdf

import "github.com/signintech/gopdf"

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func (this *Rect) ToGoRect() gopdf.Rect {
	return gopdf.Rect{W: this.W, H: this.H}
}

func (this *Rect) FromArray(rect *[4]float64) {
	this.X = rect[0]
	this.Y = rect[1]
	this.W = rect[2]
	this.H = rect[3]
}
