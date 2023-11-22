package lnpdf

type Style struct {
	FontFamily          string
	FontSize            float64
	Color               string
	Underline           bool
	HorizontalAlignment StyleHorizontalAlignment
	VerticalAlignment   StyleVerticalAlignment
	Overflow            StyleOverflow
	Multiline           bool
}

type StyleOverflow string

const (
	Crop      = StyleOverflow("crop")
	Overflow  = StyleOverflow("overflow")
	Autoscale = StyleOverflow("autoscale")
)

type StyleHorizontalAlignment string

const (
	Left   = StyleHorizontalAlignment("left")
	Center = StyleHorizontalAlignment("center")
	Right  = StyleHorizontalAlignment("right")
)

type StyleVerticalAlignment string

const (
	Top    = StyleVerticalAlignment("top")
	Middle = StyleVerticalAlignment("middle")
	Bottom = StyleVerticalAlignment("bottom")
)
