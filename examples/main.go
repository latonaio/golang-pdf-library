package main

import (
	"flag"

	lnpdf "github.com/latonaio/golang-pdf-library"
)

func main() {
	// get paremeters
	var isDrawGrid bool
	var isFillBackground bool
	var isFillSampleText bool
	var templateName string
	var dataSourceName string
	var outputPdfName string
	var privateKeyName string
	var certificateName string
	var chainName string
	flag.BoolVar(&isDrawGrid, "g", false, "draw grid")
	flag.BoolVar(&isFillBackground, "b", false, "fill background")
	flag.BoolVar(&isFillSampleText, "s", false, "fill sample text")
	flag.StringVar(&templateName, "t", "./examples/inputs/template.json", "template json file name")
	flag.StringVar(&dataSourceName, "d", "./examples/inputs/data.json", "data source json file name")
	flag.StringVar(&privateKeyName, "p", "", "path of private key for signing")
	flag.StringVar(&certificateName, "c", "", "path of certification for signing")
	flag.StringVar(&chainName, "h", "", "path of chain for signing")
	flag.StringVar(&outputPdfName, "o", "./examples/outputs/sample.pdf", "generate pdf file name")
	flag.Parse()

	// parameters for signing are needed both
	if (privateKeyName != "" && certificateName == "") || (privateKeyName == "" && certificateName != "") {
		flag.Usage()
	}

	// set resource path
	templatePath := lnpdf.ToPath(&templateName)
	dataSourcePath := lnpdf.ToPath(&dataSourceName)
	outputPdfPath := lnpdf.ToPath(&outputPdfName)
	var privateKeyPath string
	if privateKeyName != "" {
		privateKeyPath = lnpdf.ToPath(&privateKeyName)
	}
	var certificatePath string
	if certificateName != "" {
		certificatePath = lnpdf.ToPath(&certificateName)
	}
	var chainPath string
	if chainName != "" {
		chainPath = lnpdf.ToPath(&chainName)
	}

	// build
	lnpdf.Builder{
		DrawGrid:       isDrawGrid,
		FillBackground: isFillBackground,
		FillSampleText: isFillBackground,
		TemplatePath:   templatePath,
		DataSourcePath: dataSourcePath,
	}.Build().Output(&outputPdfPath)

	// sign
	if privateKeyPath != "" && certificatePath != "" {
		signingInfo := lnpdf.SigningInfo{
			Name:        "# Name #",
			Location:    "# Location #",
			Reason:      "# Reason #",
			ContactInfo: "# ContactInfo #",
			TsrUrl:      "https://freetsa.org/tsr",
		}
		lnpdf.Sign(&outputPdfPath, &privateKeyPath, &certificatePath, &chainPath, &signingInfo)
	}
}
