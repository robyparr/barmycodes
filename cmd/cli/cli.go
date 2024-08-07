package main

import (
	"flag"
	"strconv"
	"strings"

	"github.com/robyparr/barmycodes/internal"
)

type cli struct {
	values      []string
	barcodeType string
	fileType    string
	pdfPageSize internal.PDFPageSize
}

func parseCLI() cli {
	fileTypeFlag := flag.String("f", "png", "The output file type: pdf or png")
	barcodeTypeFlag := flag.String("t", "code128", "The barcode type: code128 or qr")
	pdfPageSizeFlag := flag.String("s", "", "PDF page size: NNxNNmm or NNxNNin")
	flag.Parse()

	return cli{
		values:      flag.Args(),
		barcodeType: strings.ToLower(*barcodeTypeFlag),
		fileType:    strings.ToLower(*fileTypeFlag),
		pdfPageSize: parsePdfPageSize(strings.ToLower(*pdfPageSizeFlag)),
	}
}

func parsePdfPageSize(str string) internal.PDFPageSize {
	if str == "" {
		return internal.PDFPageSize{}
	}

	values := strings.SplitN(str, "x", 2)
	unit := values[1][len(values[1])-2:]

	width, _ := strconv.Atoi(values[0])
	height, _ := strconv.Atoi(values[1][:len(values[1])-2])
	return internal.PDFPageSize{
		Width:  width,
		Height: height,
		Unit:   unit,
	}
}
