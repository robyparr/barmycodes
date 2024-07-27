package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	cli := cli{}
	cli.parse()

	barcodes := generateBarcodes(cli.values)

	if cli.fileType == "pdf" {
		err := savePDF(barcodes, cli.pdfPageSize)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		err := savePNGImages(barcodes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func generateBarcodes(values []string) []barcode {
	var barcodes []barcode
	for _, value := range values {
		barcode, err := newCode128BarCode(value)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating barcode:", err)
			os.Exit(1)
		}

		barcodes = append(barcodes, barcode)
	}

	return barcodes
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
