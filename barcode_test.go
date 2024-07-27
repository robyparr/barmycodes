package main

import (
	"bytes"
	"os"
	"testing"
)

func TestNewCode128BarCode(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		barcodeType string
		fixturePath string
	}{
		{"Basic barcode", "Test", "code128", "testdata/barcode_Test.png"},
		{"auto-resizing", "A long barcode with auto-resizing", "code128", "testdata/barcode_Long.png"},
		{"QR Code", "Test", "qr", "testdata/barcode_QR_Test.png"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wantContent, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			got, err := generateBarcode(testCase.value, testCase.barcodeType)

			if err != nil {
				t.Fatalf("Expected no error but got %s\n", err)
			}

			if !bytes.Equal(got.pngData, wantContent) {
				t.Errorf("Generated barcode does not match %s\n", testCase.fixturePath)
			}
		})
	}
}
