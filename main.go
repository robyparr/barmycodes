package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type cli struct {
	content     string
	fileType    string
	pdfPageSize string
}

func (c *cli) parse() {
	fileTypeFlag := flag.String("f", "png", "The output file type: pdf or png")
	pdfPageSizeFlag := flag.String("s", "", "PDF page size: NNxNNmm or NNxNNin")
	flag.Parse()

	c.content = flag.Arg(0)
	c.fileType = strings.ToLower(*fileTypeFlag)
	c.pdfPageSize = strings.ToLower(*pdfPageSizeFlag)
}

type pdfPageSize struct {
	width  int
	height int
	unit   string
}

func parsePdfPageSize(str string) pdfPageSize {
	if str == "" {
		return pdfPageSize{}
	}

	values := strings.SplitN(str, "x", 2)
	unit := values[1][len(values[1])-2:]

	fmt.Println(values)
	width, _ := strconv.Atoi(values[0])
	height, _ := strconv.Atoi(values[1][:len(values[1])-2])
	return pdfPageSize{
		width:  width,
		height: height,
		unit:   unit,
	}
}

func main() {
	cli := cli{}
	cli.parse()

	fileName := "barcode." + cli.fileType
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	pageSize := parsePdfPageSize(cli.pdfPageSize)
	err = newCode128BarCode(file, cli.fileType, cli.content, pageSize, time.Now)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating barcode:", err)
		file.Close()
		os.Exit(1)
	}

	fmt.Println("Created", fileName, "with content:", cli.content)
}
