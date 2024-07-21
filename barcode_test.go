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
		fixturePath string
	}{
		{"Basic barcode", "Test", "png", "testdata/barcode_Test.png"},
		{"auto-resizing", "A long barcode with auto-resizing", "png", "testdata/barcode_Long.png"},
		{"PDF", "Test", "pdf", "testdata/barcode_Test.pdf"},
	}

	nowFunc := func() time.Time { return time.Date(2024, 7, 21, 9, 0, 0, 0, time.UTC) }

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			want, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			got := new(bytes.Buffer)
			err = newCode128BarCode(got, testCase.fileType, testCase.content, nowFunc)

			if err != nil {
				t.Fatalf("Expected no error but got %s\n", err)
			}

			if !bytes.Equal(got.Bytes(), want) {
				t.Errorf("Generated barcode does not match %s\n", testCase.fixturePath)
			}
		})
	}
}
