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
		content     string
		pageSize    pdfPageSize
		fixturePath string
	}{
		{"PDF", "Test", pdfPageSize{}, "testdata/barcode_Test.pdf"},
		{"PDF sized mm", "Test", pdfPageSize{100, 100, "mm"}, "testdata/barcode_Test_100x100mm.pdf"},
		{"PDF sized in", "Test", pdfPageSize{10, 10, "in"}, "testdata/barcode_Test_10x10in.pdf"},
	}

	nowFunc := func() time.Time { return time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC) }
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			want, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			barcode, err := newCode128BarCode(testCase.content)
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
