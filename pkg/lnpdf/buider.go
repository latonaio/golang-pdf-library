package lnpdf

import (
	"encoding/json"
	model "latona-pdf/pkg/lnpdf/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/signintech/gopdf"
)

type Builder struct {
	DrawGrid       bool
	FillBackground bool
	FillSampleText bool
	TemplatePath   string
	DataSourcePath string
}

func (this Builder) Build() *Pdf {
	// load template json
	template := model.Pdf{}
	template.FromFile(&this.TemplatePath)

	// load data json
	var data interface{}
	{
		file, err := os.ReadFile(this.DataSourcePath)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(file, &data); err != nil {
			panic(err)
		}
	}

	// define directory
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	templateDir := filepath.Dir(this.TemplatePath)

	// build pdf
	pdf := Pdf{}

	// set configurations and start pdf
	paperSize := pdf.GetPaperSize(&template.Size, &template.Orientation)
	protections := gopdf.PDFProtectionConfig{
		/*UseProtection: true, */ // TODO: protectionを適用するとpdfのテンプレが使えない
		Permissions: gopdf.PermissionsPrint | gopdf.PermissionsCopy,
	}
	pdf.Start(gopdf.Config{PageSize: *paperSize, Protection: protections})

	// prepare fonts
	for n, f := range template.Fonts {
		var path string
		if strings.HasPrefix(f, "./") {
			path = currentDir + f[1:]
		} else if strings.HasPrefix(f, "/") {
			path = f
		} else {
			path = templateDir + "/" + f
		}

		var err error
		err = pdf.AddFont(n, path)
		if err != nil {
			panic(err)
		}
	}

	// build styles
	pdf.AddStyles(&template.Styles)

	// add initial page
	pdf.AddPage()

	// import pdf
	{
		var path string
		if strings.HasPrefix(template.Image, "./") {
			path = currentDir + template.Image[1:]
		} else if strings.HasPrefix(template.Image, "/") {
			path = template.Image
		} else {
			path = templateDir + "/" + template.Image
		}

		tpl := pdf.ImportPage(path, 1, "/MediaBox")
		pdf.UseImportedTemplate(tpl, 0, 0, 841.89, 595.28 /*paperSize.W, paperSize.H*/)
	}

	// draw grid
	if this.DrawGrid {
		pdf.SetTransparency(gopdf.Transparency{Alpha: 1, BlendModeType: gopdf.Color})
		pdf.SetStrokeColor(255, 0, 0)
		pdf.SetLineWidth(0.1)

		for x := 0; x < int(paperSize.W)/10; x++ {
			if x%5 == 0 {
				pdf.SetLineType("solid")
			} else {
				pdf.SetLineType("dotted")
			}
			pdf.Line(float64(x*10), 0, float64(x*10), paperSize.H)
		}

		for y := 0; y < int(paperSize.H)/10; y++ {
			if y%5 == 0 {
				pdf.SetLineType("solid")
			} else {
				pdf.SetLineType("dotted")
			}
			pdf.Line(0, float64(y*10), paperSize.W, float64(y*10))
		}
	}

	// fill background
	if this.FillBackground {
		pdf.SetTransparency(gopdf.Transparency{Alpha: 0.2, BlendModeType: gopdf.Color})
		pdf.SetFillColor(0, 0, 0)
		pdf.SetStrokeColor(0, 0, 0)
		for _, v := range template.Fields {
			pdf.Rectangle(v.Rect[0], v.Rect[1], v.Rect[0]+v.Rect[2], v.Rect[1]+v.Rect[3], "FD", 0, 0)
		}
	}

	// draw text
	{
		var drawText func(*Rect, *string, *string)
		if this.FillSampleText {
			drawText = func(rect *Rect, _ *string, styleName *string) {
				pdf.FillText(rect, styleName)
			}
		} else {
			drawText = func(rect *Rect, dataSource *string, styleName *string) {
				val := data.(map[string]interface{})[*dataSource]
				if val != nil {
					value := val.(string)
					pdf.DrawText(rect, styleName, &value)
				}
			}
		}

		rect := Rect{}
		for _, v := range template.Fields {
			if len(v.Record.Fields) > 0 {
				continue
			}

			rect.FromArray(&v.Rect)
			drawText(&rect, &v.DataSource, &v.Style)
		}
	}

	// fill sample record text
	{
		// director
		nextRecord := func(offset *gopdf.Point, field *model.Field) bool {
			if field.Record.Direction == "x" {
				offset.X += field.Record.Size[0]
				if field.Rect[2] >= offset.X+field.Record.Size[0] {
					return true
				}

				offset.X = 0
				offset.Y += field.Record.Size[1]
				if field.Rect[3] >= offset.Y+field.Record.Size[1] {
					return true
				}

				return false
			} else {
				offset.Y += field.Record.Size[1]
				if field.Rect[3] >= offset.Y+field.Record.Size[1] {
					return true
				}

				offset.Y = 0
				offset.X += field.Record.Size[0]
				if field.Rect[2] >= offset.X+field.Record.Size[0] {
					return true
				}

				return false

			}
		}

		if this.FillSampleText {
			for _, r := range template.Fields {
				if len(r.Record.Fields) == 0 {
					continue
				}

				offset := gopdf.Point{}

				// show 100 records
				for i := 0; i < 100; i++ {
					// draw all fields in record
					rect := Rect{}
					for _, f := range r.Record.Fields {
						rect.FromArray(&([4]float64{
							r.Rect[0] + offset.X + f.Rect[0],
							r.Rect[1] + offset.Y + f.Rect[1],
							f.Rect[2],
							f.Rect[3]}))
						pdf.FillText(&rect, &f.Style)
					}

					if nextRecord(&offset, &r) {
						continue
					} else {
						break
					}
				}
			}
		} else {
			for _, r := range template.Fields {
				if len(r.Record.Fields) == 0 {
					continue
				}

				records := data.(map[string]interface{})[r.DataSource]
				if records == nil {
					continue
				}

				offset := gopdf.Point{}

				for _, d := range records.([]interface{}) {
					// draw all fields in record
					rect := Rect{}
					for _, f := range r.Record.Fields {
						rect.FromArray(&([4]float64{
							r.Rect[0] + offset.X + f.Rect[0],
							r.Rect[1] + offset.Y + f.Rect[1],
							f.Rect[2],
							f.Rect[3]}))

						val := d.(map[string]interface{})[f.DataSource]
						if val != nil {
							value := val.(string)
							pdf.DrawText(&rect, &f.Style, &value)
						}
					}

					if nextRecord(&offset, &r) {
						continue
					} else {
						break
					}
				}
			}

		}

	}

	return &pdf
}
