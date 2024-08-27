package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/robyparr/barmycodes/internal"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed assets/*
var assetsFS embed.FS
var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))

type Router struct {
	NowFunc func() time.Time
	http.Handler
}

type nowFunc func() time.Time

type viewModel struct {
	Barcodes     []internal.Barcode
	BarcodeType  string
	ErrorMessage string
}

func (router Router) mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "404: The page you're looking for could not be found.")
		return
	}

	query := r.URL.Query()
	vm := viewModel{BarcodeType: query.Get("type")}
	barcodes, err := internal.GenerateBarcodes(query["b[]"], vm.BarcodeType)
	if err != nil {
		vm.ErrorMessage = "An unexpected error occurred while generating barcodes."
	} else {
		vm.Barcodes = barcodes
	}

	tmpl.ExecuteTemplate(w, "index.html.tmpl", vm)
}

func (router Router) downloadPNGHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	barcode, err := internal.GenerateBarcode(query.Get("b[]"), query.Get("type"))

	if err != nil {
		vm := viewModel{ErrorMessage: "An unexpected error occurred while generating the barcode."}
		tmpl.ExecuteTemplate(w, "index.html.tmpl", vm)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=barmycodes.png")
	w.Write(barcode.PngData)
}

func (router Router) downloadPDFHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	vm := viewModel{BarcodeType: query.Get("type")}
	barcodes, err := internal.GenerateBarcodes(query["b[]"], vm.BarcodeType)
	if err != nil {
		vm.ErrorMessage = "An unexpected error occurred while generating barcodes."
	} else {
		vm.Barcodes = barcodes
	}

	pageSize, err := getPDFPageSize(&query)
	if err != nil {
		vm.ErrorMessage = "Error parsing PDF width and height."
		tmpl.ExecuteTemplate(w, "index.html.tmpl", vm)
		return
	}

	pdf := internal.NewPdf(*pageSize, router.NowFunc)
	for _, barcode := range vm.Barcodes {
		pdf.AddBarcode(barcode)
	}

	buffer := new(bytes.Buffer)
	pdf.Write(buffer)

	w.Header().Set("Content-Disposition", "attachment; filename=barmycodes.pdf")
	w.Write(buffer.Bytes())
}

func NewRouter(now nowFunc) Router {
	router := Router{NowFunc: now}

	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.FileServerFS(assetsFS))
	mux.HandleFunc("GET /png", router.downloadPNGHandler)
	mux.HandleFunc("GET /pdf", router.downloadPDFHandler)
	mux.HandleFunc("GET /", router.mainHandler)

	router.Handler = mux
	return router
}

func getPDFPageSize(query *url.Values) (*internal.PDFPageSize, error) {
	pageSize := internal.PDFPageSize{}
	if unit := query.Get("measurement"); unit != "auto" && unit != "" {
		pageSize.Unit = unit

		width, err := strconv.Atoi(query.Get("width"))
		if err != nil {
			return nil, err
		}
		pageSize.Width = width

		height, err := strconv.Atoi(query.Get("height"))
		if err != nil {
			return nil, err
		}
		pageSize.Height = height
	}

	return &pageSize, nil
}
