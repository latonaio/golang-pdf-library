package main

import (
	"flag"
	lncommon "pdf/components/common"
	lnpdf "pdf/components/pdf"
	"pdf/components/signer"
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
	flag.StringVar(&templateName, "t", "./resources/sampleTemplate.json", "template json file name")
	flag.StringVar(&dataSourceName, "d", "./resources/sampleData.json", "data source json file name")
	flag.StringVar(&privateKeyName, "p", "", "path of private key for signing")
	flag.StringVar(&certificateName, "c", "", "path of certification for signing")
	flag.StringVar(&chainName, "h", "", "path of chain for signing")
	flag.StringVar(&outputPdfName, "o", "./sample.pdf", "generate pdf file name")
	flag.Parse()

	// parameters for signing are needed both
	if (privateKeyName != "" && certificateName == "") || (privateKeyName == "" && certificateName != "") {
		flag.Usage()
	}

	// set resource path
	templatePath := lncommon.ToPath(&templateName)
	dataSourcePath := lncommon.ToPath(&dataSourceName)
	outputPdfPath := lncommon.ToPath(&outputPdfName)
	var privateKeyPath string
	if privateKeyName != "" {
		privateKeyPath = lncommon.ToPath(&privateKeyName)
	}
	var certificatePath string
	if certificateName != "" {
		certificatePath = lncommon.ToPath(&certificateName)
	}
	var chainPath string
	if chainName != "" {
		chainPath = lncommon.ToPath(&chainName)
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
		signingInfo := signer.SigningInfo{
			Name:        "# Name #",
			Location:    "# Location #",
			Reason:      "# Reason #",
			ContactInfo: "# ContactInfo #",
			TsrUrl:      "https://freetsa.org/tsr",
		}
		signer.Sign(&outputPdfPath, &privateKeyPath, &certificatePath, &chainPath, &signingInfo)
	}
}
