package lnpdf

import (
	"math"
	lncommon "pdf/components/common"
	"pdf/components/rgb"
	"strings"

	"github.com/signintech/gopdf"
)

// gopdf.GoPdf structure with style
type Pdf struct {
	gopdf.GoPdf

	Styles *map[string]Style
	Rect   Rect
}

// add .ttf font
func (this *Pdf) AddFont(family string, path string) error {
	return this.AddTTFFont(family, path)
}

// draw text within rect with style
func (this *Pdf) DrawText(rect *Rect, styleName *string, text *string) {
	option := gopdf.CellOption{}

	// set horizontal alignment
	style := (*this.Styles)[*styleName]
	if style.HorizontalAlignment == Right {
		option.Align |= gopdf.Right
	} else if style.HorizontalAlignment == Center {
		option.Align |= gopdf.Center
	} else {
		option.Align |= gopdf.Left
	}

	// set vertical alignment
	if style.VerticalAlignment == Bottom {
		option.Align |= gopdf.Bottom
	} else if style.VerticalAlignment == Middle {
		option.Align |= gopdf.Middle
	} else {
		option.Align |= gopdf.Top
	}

	// set color
	this.SetTransparency(gopdf.Transparency{Alpha: 1, BlendModeType: gopdf.Color})
	hex := rgb.Hex(style.Color)
	rgb, err := hex.ToRGB()
	if err != nil {
		panic(err)
	}
	this.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)

	// set font style (only underline)
	var fontStyle string
	if style.Underline {
		fontStyle = "U"
	}

	// apply font style
	this.SetFont(style.FontFamily, fontStyle, style.FontSize)

	// process overflow
	adjustedText := text
	if style.Overflow == Crop {
		// crop text
		adjustedText = this.getCroppedText(rect, styleName, style.Multiline, text)
	} else if style.Overflow == Autoscale {
		// scale font size
		fontSize := this.getAdjustedFontSize(rect, styleName, style.Multiline, text)
		this.SetFont(style.FontFamily, fontStyle, fontSize)
	} else {
		// do nothing even if it overflows
		// nop
	}

	// draw text
	this.SetXY(rect.X, rect.Y)
	goRect := rect.ToGoRect()
	if style.Multiline {
		this.MultiCellWithOption(&goRect, *adjustedText, option)
	} else {
		this.CellWithOption(&goRect, *adjustedText, option)
	}
}

// draw sample text within rect
func (this *Pdf) FillText(rect *Rect, styleName *string) {
	style := (*this.Styles)[*styleName]
	this.DrawText(rect, styleName, this.getFillText(rect, styleName, style.Multiline))
}

// get sample text within rect
func (this *Pdf) getFillText(rect *Rect, styleName *string, multiline bool) *string {
	style := (*this.Styles)[*styleName]

	// set font
	var fontStyle string
	if style.Underline {
		fontStyle = "U"
	}
	this.SetFont(style.FontFamily, fontStyle, style.FontSize)

	// single line text
	const max = 100
	sample := strings.Repeat("X", max)
	width, err := this.MeasureTextWidth(sample)
	if err != nil {
		panic(err)
	}
	if width > rect.W {
		// adjust text length
		sample = strings.Repeat("X", int(math.Floor((rect.W/width)*max)))
	}

	// multi line text
	if multiline {
		height, _ := this.MeasureCellHeightByText(sample)
		if rect.H > height {
			// repeat until rect.H
			sample = strings.Repeat(sample, int(math.Floor((rect.H / height))))
		}
	}

	return &sample
}

// get sample text within rect
func (this *Pdf) getCroppedText(rect *Rect, styleName *string, multiline bool, text *string) *string {
	// AutoScale is not applicable to multiline text
	if multiline {
		return text
	}

	style := (*this.Styles)[*styleName]

	// set font
	var fontStyle string
	if style.Underline {
		fontStyle = "U"
	}
	this.SetFont(style.FontFamily, fontStyle, style.FontSize)

	// single line text
	sample := *text
	for {
		width, err := this.MeasureTextWidth(sample)
		if err != nil {
			break
		}
		if width > rect.W {
			sample = sample[:len(sample)-1]
			if sample == "" {
				break
			}
		} else {
			break
		}
	}

	return &sample
}

// get fontsize within rect
func (this *Pdf) getAdjustedFontSize(rect *Rect, styleName *string, multiline bool, text *string) float64 {
	style := (*this.Styles)[*styleName]

	// AutoScale is not applicable to multiline text
	if multiline {
		return style.FontSize
	}

	// set font
	var fontStyle string
	if style.Underline {
		fontStyle = "U"
	}

	// single line text
	fontSize := style.FontSize
	for {
		width, err := this.MeasureTextWidth(*text)
		if err != nil {
			break
		}
		if width > rect.W {
			fontSize -= 0.25
			this.SetFont(style.FontFamily, fontStyle, fontSize)
			if fontSize < 2 {
				break
			}
		} else {
			break
		}
	}

	return fontSize
}

func (this *Pdf) GetPaperSize(paperSize *string, orientation *string) *gopdf.Rect {
	var goRect gopdf.Rect
	if false {
		// nop
	} else if *paperSize == string(lncommon.A1) {
		goRect = *gopdf.PageSizeA1
	} else if *paperSize == string(lncommon.A2) {
		goRect = *gopdf.PageSizeA2
	} else if *paperSize == string(lncommon.A3) {
		goRect = *gopdf.PageSizeA3
	} else if *paperSize == string(lncommon.A5) {
		goRect = *gopdf.PageSizeA5
	} else if *paperSize == string(lncommon.B4) {
		goRect = *gopdf.PageSizeB4
	} else if *paperSize == string(lncommon.B5) {
		goRect = *gopdf.PageSizeB5
	} else {
		goRect = *gopdf.PageSizeA4
	}

	// portlait -> landscape
	if *orientation == string(lncommon.Landscape) {
		goRect = gopdf.Rect{W: goRect.H, H: goRect.W}
	}

	return &goRect
}

func (this *Pdf) Output(path *string) {
	this.WritePdf(*path)
}
