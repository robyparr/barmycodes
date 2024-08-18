package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
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

func (router Router) mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "404: The page you're looking for could not be found.")
		return
	}

	query := r.URL.Query()
	barcodes, _ := internal.GenerateBarcodes(query["b[]"], query.Get("type"))

	tmpl.ExecuteTemplate(w, "index.html.tmpl", barcodes)
}

func (router Router) downloadPNGHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	barcode, _ := internal.GenerateBarcode(query.Get("b[]"), query.Get("type"))

	w.Header().Set("Content-Disposition", "attachment; filename=barmycodes.png")
	w.Write(barcode.PngData)
}

func (router Router) downloadPDFHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	barcodes, _ := internal.GenerateBarcodes(query["b[]"], query.Get("type"))

	pageSize := internal.PDFPageSize{}
	if unit := query.Get("measurement"); unit != "auto" {
		pageSize.Unit = unit

		width, _ := strconv.Atoi(query.Get("width"))
		pageSize.Width = width

		height, _ := strconv.Atoi(query.Get("height"))
		pageSize.Height = height
	}

	buffer := new(bytes.Buffer)
	pdf := internal.NewPdf(pageSize, router.NowFunc)
	for _, barcode := range barcodes {
		pdf.AddBarcode(barcode)
	}
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
