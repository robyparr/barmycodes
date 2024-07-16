package main

import (
	"fmt"
	"os"
)

func main() {
	text := os.Args[1]

	file, err := os.Create("barcode.png")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	err = newCode128BarCode(file, text)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating barcode:", err)
		file.Close()
		os.Exit(1)
	}

	fmt.Println("Created barcode.png with content:", text)
}
