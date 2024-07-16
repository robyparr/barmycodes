package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	labelFontFace  = basicfont.Face7x13
	labelMarginTop = 2

	barcodePadding = 10
)

func newCode128BarCode(w io.Writer, text string) error {
	bcode, err := code128.Encode(text)
	if err != nil {
		return err
	}

	scaledBc, err := barcode.Scale(bcode, bcode.Bounds().Dx()*2, 150)
	if err != nil {
		return err
	}

	err = png.Encode(w, addLabel(scaledBc))
	if err != nil {
		return err
	}

	return nil
}

func addLabel(bcode barcode.Barcode) image.Image {
	// Based on https://github.com/boombuler/barcode/wiki/Content-String
	label := bcode.Content()
	labelBounds, _ := font.BoundString(labelFontFace, label)
	labelWidth := int((labelBounds.Max.X - labelBounds.Min.X) / 64)
	labelHeight := int((labelBounds.Max.Y - labelBounds.Min.Y) / 64)

	imgHeight := labelHeight + bcode.Bounds().Dy() + labelMarginTop
	imgWidth := labelWidth
	if bcode.Bounds().Dx() > imgWidth {
		imgWidth = bcode.Bounds().Dx()
	}

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	backgroundRect := image.Rect(0, 0, imgWidth, imgHeight)
	draw.Draw(img, backgroundRect, &image.Uniform{color.White}, bcode.Bounds().Min, draw.Over)

	barcodeRect := image.Rect(0, 0, bcode.Bounds().Dx(), bcode.Bounds().Dy())
	draw.Draw(img, barcodeRect, bcode, bcode.Bounds().Min, draw.Over)

	labelOffsetY := bcode.Bounds().Dy() + labelMarginTop - int(labelBounds.Min.Y/64)
	labelOffsetX := (imgWidth - labelWidth) / 2

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
