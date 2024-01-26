package model

type Pdf struct {
	Version     string
	Orientation string
	Size        string
	Image       string
	Fonts       map[string]string
	Styles      map[string]Style
	Fields      []Field
}
