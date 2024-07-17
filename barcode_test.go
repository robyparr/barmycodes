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
			want, err := os.ReadFile(testCase.fixturePath)
			if err != nil {
				t.Fatal(err)
			}

			got := new(bytes.Buffer)
			newCode128BarCode(got, testCase.content)

			if !bytes.Equal(got.Bytes(), want) {
				t.Errorf("Generated barcode does not match %s\n", testCase.fixturePath)
			}
		})
	}
}
