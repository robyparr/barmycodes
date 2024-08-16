package internal

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"

	bBarcode "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	_ "embed"
)

//go:embed fonts/opensans/OpenSans-Regular.ttf
var fontFile []byte

const (
	labelMarginTop = 10
	imgPaddingX    = 20
	imgPaddingY    = 5
	qrCodePadding  = 60
)

type Barcode struct {
	PngData []byte
	Type    string
	Value   string
	width   int
	height  int
}

func (b Barcode) Base64() string {
	return base64.RawStdEncoding.EncodeToString(b.PngData)
}

func GenerateBarcode(text string, barcodeType string) (Barcode, error) {
	bc := Barcode{Type: barcodeType, Value: text}
	var img image.Image
	var err error

	if strings.ToLower(barcodeType) == "qr" {
		img, err = generateQRCode(&bc)
	} else {
		img, err = generateCode128Barcode(&bc)
	}

	pngBuffer := new(bytes.Buffer)
	err = png.Encode(pngBuffer, img)
	if err != nil {
		return Barcode{}, err
	}

	bc.PngData = pngBuffer.Bytes()
	return bc, nil
}

func GenerateBarcodes(values []string, barcodeType string) ([]Barcode, error) {
	var barcodes []Barcode
	for _, value := range values {
		barcode, err := GenerateBarcode(value, barcodeType)
		if err != nil {
			return barcodes, err
		}

		barcodes = append(barcodes, barcode)
	}

	return barcodes, nil
}

func generateCode128Barcode(bc *Barcode) (image.Image, error) {
	bcode, err := code128.Encode(bc.Value)
	if err != nil {
		return nil, err
	}

	bc.height = 120
	bc.width = bcode.Bounds().Dx() * 3
	if len(bc.Value) >= 10 {
		bc.width += len(bc.Value) * 15
	}

	scaledBc, err := bBarcode.Scale(bcode, bc.width, bc.height)
	if err != nil {
		return nil, err
	}

	return addLabel(scaledBc), nil
}

func generateQRCode(bc *Barcode) (image.Image, error) {
	bc.height = 245
	bc.width = 245
	bcode, err := qr.Encode(bc.Value, qr.H, qr.Unicode)
	if err != nil {
		return nil, err
	}

	scaledBc, err := bBarcode.Scale(bcode, bc.width, bc.height)
	if err != nil {
		return nil, err
	}

	return addBorder(scaledBc), nil
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

func addBorder(bcode bBarcode.Barcode) image.Image {
	imgRect := image.Rect(0, 0, bcode.Bounds().Dx()+qrCodePadding, bcode.Bounds().Dy()+qrCodePadding)
	img := image.NewRGBA(imgRect)
	draw.Draw(img, imgRect, &image.Uniform{color.White}, bcode.Bounds().Min, draw.Over)

	offset := qrCodePadding / 2
	barcodeRect := image.Rect(offset, offset, bcode.Bounds().Dx()+offset, bcode.Bounds().Dy()+offset)
	draw.Draw(img, barcodeRect, bcode, bcode.Bounds().Min, draw.Over)

	return img
}
