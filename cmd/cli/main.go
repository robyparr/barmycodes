package main

import (
	"fmt"
	"os"
	"time"

	"github.com/robyparr/barmycodes/internal"
)

func main() {
	cli := parseCLI()

	barcodes, err := internal.GenerateBarcodes(cli.values, cli.barcodeType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if cli.fileType == "pdf" {
		err = savePDF(barcodes, cli.pdfPageSize)
	} else {
		err = savePNGImages(barcodes)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func savePDF(barcodes []internal.Barcode, pageSize internal.PDFPageSize) error {
	pdf := internal.NewPdf(pageSize, time.Now)
	for _, barcode := range barcodes {
		pdf.AddBarcode(barcode)
	}

	filename := "barmycodes.pdf"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}

	err = pdf.Write(file)
	if err != nil {
		file.Close()
		os.Remove(filename)
		return fmt.Errorf("Error writing PDF: %s", err)
	}

	fmt.Println("Created", filename, "with", len(barcodes), "barcodes")
	return nil
}

func savePNGImages(barcodes []internal.Barcode) error {
	for _, barcode := range barcodes {
		filename := fmt.Sprintf("%s.png", barcode.Value)

		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("Error creating file: %s", err)
		}
		defer file.Close()

		_, err = file.Write(barcode.PngData)
		if err != nil {
			file.Close()
			os.Remove(filename)
			return fmt.Errorf("Error writing file: %s", err)
		}

		fmt.Println("Created", filename, "with content:", barcode.Value)
	}

	return nil
}
