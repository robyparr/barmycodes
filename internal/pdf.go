package internal

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

type PDFPageSize struct {
	Width  int
	Height int
	Unit   string
}

func (p PDFPageSize) isZero() bool {
	return p.Unit == ""
}

type pdf struct {
	fpdf     *fpdf.Fpdf
	pageSize fpdf.SizeType
}

func (p *pdf) AddBarcode(bc Barcode) {
	pageSize := fpdf.SizeType{Wd: float64(bc.width), Ht: float64(bc.height + 20)}
	if p.pageSize.Wd != 0 {
		pageSize.Wd = p.pageSize.Wd
		pageSize.Ht = p.pageSize.Ht
	}

	pngDataReader := bytes.NewReader(bc.PngData)
	p.fpdf.AddPageFormat(pageFormatPortrait, pageSize)
	p.fpdf.RegisterImageOptionsReader(bc.Value, fpdf.ImageOptions{ImageType: "PNG"}, pngDataReader)
	p.fpdf.ImageOptions(bc.Value, 0, 0, pageSize.Wd, pageSize.Ht, false, fpdf.ImageOptions{}, 0, "")
}

func (p *pdf) Write(w io.Writer) error {
	err := p.fpdf.Output(w)
	if err != nil {
		return err
	}

	return nil
}

func NewPdf(pageSize PDFPageSize, now nowFunc) pdf {
	pdf := pdf{}
	unit := defaultUnit
	if !pageSize.isZero() {
		unit = pageSize.Unit
		pdf.pageSize.Wd = float64(pageSize.Width)
		pdf.pageSize.Ht = float64(pageSize.Height)
	}

	pdf.fpdf = fpdf.New(pageFormatPortrait, unit, pageSizeA4, "")
	pdf.fpdf.SetCreationDate(now().UTC())
	pdf.fpdf.SetModificationDate(now().UTC())

	return pdf
}
