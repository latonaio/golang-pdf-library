package lnpdf

import model "latona-pdf/pkg/lnpdf/models"

/*
	FontFamily          string
	FontStyle           string
	FontSize            float64
	AutoScale           bool
	ForeColor           string
	BackColor           string
	Underline           bool
	HorizontalAlignment string
	VerticalAlignment   string
	Format              string
	IsMultiline         bool
	IsOverflow          bool
	IsReduce            bool

*/

func (this *Pdf) AddStyles(styles *map[string]model.Style) {
	(*this).Styles = &map[string]Style{}

	// 全ての model.styles を pdf.styles に追加
	for n, v := range *styles {
		route := []string{n}
		route = *getStyleRoute(&route, &v, styles)

		newStyle := Style{}
		for _, name := range route {
			s := (*styles)[name]
			newStyle.mergeStyle(&s)
		}

		(*this.Styles)[n] = newStyle
	}
}

func getStyleRoute(route *[]string, style *model.Style, styles *map[string]model.Style) *[]string {
	if style.Style == nil {
		return route
	}

	baseStyleName := *style.Style
	baseStyle := (*styles)[baseStyleName]

	newRoute := append([]string{baseStyleName}, *route...)

	return getStyleRoute(&newRoute, &baseStyle, styles)
}

func (this *Style) mergeStyle(m *model.Style) {
	if m.FontFamily != nil {
		this.FontFamily = *m.FontFamily
	}
	if m.FontSize != nil {
		this.FontSize = *m.FontSize
	}
	if m.Color != nil {
		this.Color = *m.Color
	}
	if m.Underline != nil {
		this.Underline = *m.Underline
	}
	if m.HorizontalAlignment != nil {
		if *m.HorizontalAlignment == string(Center) {
			this.HorizontalAlignment = Center
		} else if *m.HorizontalAlignment == string(Right) {
			this.HorizontalAlignment = Right
		} else {
			this.HorizontalAlignment = Left
		}
	}
	if m.VerticalAlignment != nil {
		if *m.VerticalAlignment == string(Top) {
			this.VerticalAlignment = Top
		} else if *m.VerticalAlignment == string(Bottom) {
			this.VerticalAlignment = Bottom
		} else {
			this.VerticalAlignment = Middle
		}
	}
	if m.Multiline != nil {
		this.Multiline = *m.Multiline
	}
	if m.Overflow != nil {
		if *m.Overflow == string(Overflow) {
			this.Overflow = Overflow
		} else if *m.Overflow == string(Autoscale) {
			this.Overflow = Autoscale
		} else {
			this.Overflow = Crop
		}
	}
}
