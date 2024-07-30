package internal_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/robyparr/barmycodes/internal"
)

func TestPDFCreation(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		barcodeType string
		pageSize    internal.PDFPageSize
		fixturePath string
	}{
		{"PDF", "Test", "code128", internal.PDFPageSize{}, "testdata/barcode_Test.pdf"},
		{"PDF sized mm", "Test", "code128", internal.PDFPageSize{100, 100, "mm"}, "testdata/barcode_Test_100x100mm.pdf"},
		{"PDF sized in", "Test", "code128", internal.PDFPageSize{10, 10, "in"}, "testdata/barcode_Test_10x10in.pdf"},
		{"QR PDF", "Test", "qr", internal.PDFPageSize{}, "testdata/barcode_QR_Test.pdf"},
		{"QR PDF sized mm", "Test", "qr", internal.PDFPageSize{200, 200, "mm"}, "testdata/barcode_QR_Test_200x200mm.pdf"},
	}

	nowFunc := func() time.Time { return time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC) }
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			want, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			barcode, err := internal.GenerateBarcode(testCase.value, testCase.barcodeType)
			if err != nil {
				t.Fatal(err)
			}

			pdf := internal.NewPdf(testCase.pageSize, nowFunc)
			pdf.AddBarcode(barcode)

			got := new(bytes.Buffer)
			err = pdf.Write(got)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got.Bytes(), want) {
				t.Errorf("Generated PDF does not match %s", testCase.fixturePath)
			}
		})
	}
}
