package main

import (
	"bytes"
	"io"
	"time"

	"github.com/go-pdf/fpdf"
)

const (
	defaultUnit        string = "pt"
	pageFormatPortrait string = "P"
	pageSizeA4         string = "A4"
)

type nowFunc = func() time.Time

type pdfPageSize struct {
	width  int
	height int
	unit   string
}

func (p pdfPageSize) isZero() bool {
	return p.unit == ""
}

type pdf struct {
	fpdf     *fpdf.Fpdf
	pageSize fpdf.SizeType
}

func (p *pdf) addBarcode(bc barcode) {
	pageSize := fpdf.SizeType{Wd: float64(bc.width), Ht: float64(bc.height + 20)}
	if p.pageSize.Wd != 0 {
		pageSize.Wd = p.pageSize.Wd
		pageSize.Ht = p.pageSize.Ht
	}

	pngDataReader := bytes.NewReader(bc.pngData)
	p.fpdf.AddPageFormat(pageFormatPortrait, pageSize)
	p.fpdf.RegisterImageOptionsReader(bc.value, fpdf.ImageOptions{ImageType: "PNG"}, pngDataReader)
	p.fpdf.ImageOptions(bc.value, 0, 0, pageSize.Wd, pageSize.Ht, false, fpdf.ImageOptions{}, 0, "")
}

func newPdf(pageSize pdfPageSize, now nowFunc) pdf {
	pdf := pdf{}
	unit := defaultUnit
	if !pageSize.isZero() {
		unit = pageSize.unit
		pdf.pageSize.Wd = float64(pageSize.width)
		pdf.pageSize.Ht = float64(pageSize.height)
	}

	pdf.fpdf = fpdf.New(pageFormatPortrait, unit, pageSizeA4, "")
	pdf.fpdf.SetCreationDate(now())
	pdf.fpdf.SetModificationDate(now())

	return pdf
}

func (p *pdf) write(w io.Writer) error {
	err := p.fpdf.Output(w)
	if err != nil {
		return err
	}

	return nil
}
