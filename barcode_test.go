package main

import (
	"bytes"
	"os"
	"testing"
)

func TestNewCode128BarCode(t *testing.T) {
	want, err := os.ReadFile("testdata/barcode_Test.png")
	if err != nil {
		t.Fatal(err)
	}

	got := new(bytes.Buffer)
	newCode128BarCode(got, "Test")

	if !bytes.Equal(got.Bytes(), want) {
		t.Error("Generated barcode does not match testdata/barcode_Test.png")
	}
}
