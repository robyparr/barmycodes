package main

import (
	"bytes"
	"os"
	"testing"
)

func TestNewCode128BarCode(t *testing.T) {
	testCases := []struct {
		name        string
		content     string
		fixturePath string
	}{
		{"Basic barcode", "Test", "testdata/barcode_Test.png"},
		{"auto-resizing", "A long barcode with auto-resizing", "testdata/barcode_Long.png"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wantContent, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			got, err := newCode128BarCode(testCase.content)

			if err != nil {
				t.Fatalf("Expected no error but got %s\n", err)
			}

			if !bytes.Equal(got.pngData, wantContent) {
				t.Errorf("Generated barcode does not match %s\n", testCase.fixturePath)
			}
		})
	}
}
