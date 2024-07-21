package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type cli struct {
	content  string
	fileType string
}

func (c *cli) parse() {
	fileTypeFlag := flag.String("f", "png", "The output file type: pdf or png")
	flag.Parse()

	c.content = flag.Arg(0)
	c.fileType = strings.ToLower(*fileTypeFlag)
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

	err = newCode128BarCode(file, cli.fileType, cli.content, time.Now)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating barcode:", err)
		file.Close()
		os.Exit(1)
	}

	fmt.Println("Created", fileName, "with content:", cli.content)
}
