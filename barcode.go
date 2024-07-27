package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	bBarcode "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	_ "embed"
)

var (
	//go:embed fonts/opensans/OpenSans-Regular.ttf
	fontFile       []byte
	labelMarginTop = 10
	imgPaddingX    = 20
	imgPaddingY    = 5
)

type barcode struct {
	pngData []byte
	value   string
	width   int
	height  int
}

func newCode128BarCode(text string) (barcode, error) {
	bc := barcode{value: text, height: 120}
	bcode, err := code128.Encode(text)
	if err != nil {
		return barcode{}, err
	}

	bc.width = bcode.Bounds().Dx() * 3
	textLen := len(text)
	if textLen >= 10 {
		bc.width += textLen * 15
	}

	scaledBc, err := bBarcode.Scale(bcode, bc.width, bc.height)
	if err != nil {
		return barcode{}, err
	}

	pngBuffer := new(bytes.Buffer)
	err = png.Encode(pngBuffer, addLabel(scaledBc))
	if err != nil {
		return barcode{}, err
	}

	bc.pngData = pngBuffer.Bytes()
	return bc, nil
}

func addLabel(bcode bBarcode.Barcode) image.Image {
	labelFont, _ := opentype.Parse(fontFile)
	labelFontFace, _ := opentype.NewFace(labelFont, &opentype.FaceOptions{
		Size:    float64(20),
		DPI:     200,
		Hinting: font.HintingNone,
	})

	// Based on https://github.com/boombuler/barcode/wiki/Content-String
	label := bcode.Content()
	labelBounds, _ := font.BoundString(labelFontFace, label)
	labelWidth := int((labelBounds.Max.X - labelBounds.Min.X) / 64)
	labelHeight := int((labelBounds.Max.Y - labelBounds.Min.Y) / 64)

	imgHeight := labelHeight + bcode.Bounds().Dy() + labelMarginTop + (imgPaddingY * 2)
	imgWidth := labelWidth
	if bcode.Bounds().Dx() > imgWidth {
		imgWidth = bcode.Bounds().Dx()
	}
	imgWidth += (imgPaddingX * 2)

	imgRect := image.Rect(-(imgPaddingX * 2), -(imgPaddingY * 2), imgWidth, imgHeight)
	img := image.NewRGBA(imgRect)
	draw.Draw(img, imgRect, &image.Uniform{color.White}, bcode.Bounds().Min, draw.Over)

	barcodeRect := image.Rect(0, 0, bcode.Bounds().Dx(), bcode.Bounds().Dy())
	draw.Draw(img, barcodeRect, bcode, bcode.Bounds().Min, draw.Over)

	labelOffsetY := bcode.Bounds().Dy() + labelMarginTop - int(labelBounds.Min.Y/64)
	labelOffsetX := ((imgWidth - labelWidth) / 2) - imgPaddingX

	point := fixed.Point26_6{
		X: fixed.Int26_6(labelOffsetX * 64),
		Y: fixed.Int26_6(labelOffsetY * 64),
	}

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: labelFontFace,
		Dot:  point,
	}
	drawer.DrawString(label)

	return img
}
