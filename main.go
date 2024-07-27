package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	cli := cli{}
	cli.parse()

	barcodes, err := generateBarcodes(cli.values, cli.barcodeType)
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

func generateBarcodes(values []string, barcodeType string) ([]barcode, error) {
	var barcodes []barcode
	for _, value := range values {
		barcode, err := generateBarcode(value, barcodeType)
		if err != nil {
			return barcodes, err
		}

		barcodes = append(barcodes, barcode)
	}

	return barcodes, nil
}

func savePDF(barcodes []barcode, pageSize pdfPageSize) error {
	pdf := newPdf(pageSize, time.Now)
	for _, barcode := range barcodes {
		pdf.addBarcode(barcode)
	}

	filename := "barmycodes.pdf"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err)
	}

	err = pdf.write(file)
	if err != nil {
		file.Close()
		os.Remove(filename)
		return fmt.Errorf("Error writing PDF: %s", err)
	}

	fmt.Println("Created", filename, "with", len(barcodes), "barcodes")
	return nil
}

func savePNGImages(barcodes []barcode) error {
	for _, barcode := range barcodes {
		filename := fmt.Sprintf("%s.png", barcode.value)

		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("Error creating file: %s", err)
		}
		defer file.Close()

		_, err = file.Write(barcode.pngData)
		if err != nil {
			file.Close()
			os.Remove(filename)
			return fmt.Errorf("Error writing file: %s", err)
		}

		fmt.Println("Created", filename, "with content:", barcode.value)
	}

	return nil
}
