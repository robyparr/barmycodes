package internal_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/robyparr/barmycodes/internal"
)

func TestGenerateBarcode(t *testing.T) {
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

			got, err := internal.GenerateBarcode(testCase.value, testCase.barcodeType)

			if err != nil {
				t.Fatalf("Expected no error but got %s\n", err)
			}

			if !bytes.Equal(got.PngData, wantContent) {
				t.Errorf("Generated barcode does not match %s\n", testCase.fixturePath)
			}
		})
	}
}
