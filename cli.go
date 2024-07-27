package main

import (
	"flag"
	"strconv"
	"strings"
)

type cli struct {
	values      []string
	barcodeType string
	fileType    string
	pdfPageSize pdfPageSize
}

func (c *cli) parse() {
	fileTypeFlag := flag.String("f", "png", "The output file type: pdf or png")
	barcodeTypeFlag := flag.String("t", "code128", "The barcode type: code128 or qr")
	pdfPageSizeFlag := flag.String("s", "", "PDF page size: NNxNNmm or NNxNNin")
	flag.Parse()

	c.values = flag.Args()
	c.barcodeType = strings.ToLower(*barcodeTypeFlag)
	c.fileType = strings.ToLower(*fileTypeFlag)
	c.pdfPageSize = parsePdfPageSize(strings.ToLower(*pdfPageSizeFlag))
}

func parsePdfPageSize(str string) pdfPageSize {
	if str == "" {
		return pdfPageSize{}
	}

	values := strings.SplitN(str, "x", 2)
	unit := values[1][len(values[1])-2:]

	width, _ := strconv.Atoi(values[0])
	height, _ := strconv.Atoi(values[1][:len(values[1])-2])
	return pdfPageSize{
		width:  width,
		height: height,
		unit:   unit,
	}
}
