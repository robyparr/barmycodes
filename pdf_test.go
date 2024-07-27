package main

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestPDFCreation(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		barcodeType string
		pageSize    pdfPageSize
		fixturePath string
	}{
		{"PDF", "Test", "code128", pdfPageSize{}, "testdata/barcode_Test.pdf"},
		{"PDF sized mm", "Test", "code128", pdfPageSize{100, 100, "mm"}, "testdata/barcode_Test_100x100mm.pdf"},
		{"PDF sized in", "Test", "code128", pdfPageSize{10, 10, "in"}, "testdata/barcode_Test_10x10in.pdf"},
		{"QR PDF", "Test", "qr", pdfPageSize{}, "testdata/barcode_QR_Test.pdf"},
		{"QR PDF sized mm", "Test", "qr", pdfPageSize{200, 200, "mm"}, "testdata/barcode_QR_Test_200x200mm.pdf"},
	}

	nowFunc := func() time.Time { return time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC) }
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			want, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			barcode, err := generateBarcode(testCase.value, testCase.barcodeType)
			if err != nil {
				t.Fatal(err)
			}

			pdf := newPdf(testCase.pageSize, nowFunc)
			pdf.addBarcode(barcode)

			got := new(bytes.Buffer)
			err = pdf.write(got)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got.Bytes(), want) {
				t.Errorf("Generated PDF does not match %s", testCase.fixturePath)
			}
		})
	}
}
