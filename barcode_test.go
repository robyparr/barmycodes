package main

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestNewCode128BarCode(t *testing.T) {
	testCases := []struct {
		name        string
		content     string
		fileType    string
		pageSize    pdfPageSize
		fixturePath string
	}{
		{"Basic barcode", "Test", "png", pdfPageSize{}, "testdata/barcode_Test.png"},
		{"auto-resizing", "A long barcode with auto-resizing", "png", pdfPageSize{}, "testdata/barcode_Long.png"},
		{"PDF", "Test", "pdf", pdfPageSize{}, "testdata/barcode_Test.pdf"},
		{"PDF sized mm", "Test", "pdf", pdfPageSize{100, 100, "mm"}, "testdata/barcode_Test_100x100mm.pdf"},
		{"PDF sized in", "Test", "pdf", pdfPageSize{10, 10, "in"}, "testdata/barcode_Test_10x10in.pdf"},
	}

	nowFunc := func() time.Time { return time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC) }

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			want, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			got := new(bytes.Buffer)
			err = newCode128BarCode(got, testCase.fileType, testCase.content, testCase.pageSize, nowFunc)

			if err != nil {
				t.Fatalf("Expected no error but got %s\n", err)
			}

			if !bytes.Equal(got.Bytes(), want) {
				t.Errorf("Generated barcode does not match %s\n", testCase.fixturePath)
			}
		})
	}
}
