package lnpdf

type PdfEntity struct {
	Version     string
	Orientation string
	Size        string
	Image       string
	Fonts       map[string]string
	Styles      map[string]StyleEntity
	Fields      []FieldEntity
}

type StyleEntity struct {
	Style               *string
	FontFamily          *string
	FontSize            *float64
	Color               *string
	Underline           *bool
	HorizontalAlignment *string
	VerticalAlignment   *string
	Format              *string
	Multiline           *bool
	Overflow            *string
}

type FieldEntity struct {
	DataSource string
	Style      string
	Rect       [4]float64
	Record     RecordEntity
}

type RecordEntity struct {
	Direction string
	Size      []float64
	Fields    []RecordFieldEntity
}

type RecordFieldEntity struct {
	DataSource string
	Style      string
	Rect       [4]float64
}
