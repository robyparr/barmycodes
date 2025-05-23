package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
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
	Barcodes               []internal.Barcode
	BarcodeType            string
	ErrorMessage           string
	EventHorizonProjectKey string
	BannerMessage          string
}

func newViewModel() viewModel {
	return viewModel{
		EventHorizonProjectKey: os.Getenv("EVENT_HORIZON_PROJECT_KEY"),
		BannerMessage:          os.Getenv("BANNER_MESSAGE"),
	}
}

func (router Router) mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "404: The page you're looking for could not be found.")
		return
	}

	query := r.URL.Query()
	vm := newViewModel()
	vm.BarcodeType = query.Get("type")
	if len(query["b[]"]) > 25 {
		vm.ErrorMessage = "You cannot generate more than 25 barcodes at one time."
		tmpl.ExecuteTemplate(w, "index.html.tmpl", vm)
		return
	}

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
		vm := newViewModel()
		vm.ErrorMessage = "An unexpected error occurred while generating the barcode."
		tmpl.ExecuteTemplate(w, "index.html.tmpl", vm)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=barmycodes.png")
	w.Write(barcode.PngData)
}

func (router Router) downloadPDFHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	vm := newViewModel()
	vm.BarcodeType = query.Get("type")

	if len(query["b[]"]) > 25 {
		vm.ErrorMessage = "You cannot generate more than 25 barcodes at one time."
		tmpl.ExecuteTemplate(w, "index.html.tmpl", vm)
		return
	}

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

func (router Router) privacyHandler(w http.ResponseWriter, r *http.Request) {
	vm := newViewModel()
	tmpl.ExecuteTemplate(w, "privacy.html.tmpl", vm)
}

func (router Router) attributionsHandler(w http.ResponseWriter, r *http.Request) {
	vm := newViewModel()
	tmpl.ExecuteTemplate(w, "attributions.html.tmpl", vm)
}

func NewRouter(now nowFunc) Router {
	router := Router{NowFunc: now}

	loggingMiddleware := func(next http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r)
			elapsedTime := time.Since(start).Truncate(time.Millisecond)

			log.Printf("%s %s in %s\n", r.Method, r.URL.Path, elapsedTime)
		})
	}

	secureHeadersMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; script-src 'self' 'unsafe-inline' https://eh.robyparr.com; style-src 'self' 'unsafe-inline'; connect-src 'self' https://eh.robyparr.com")

			w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-XSS-Protection", "0")

			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.FileServerFS(assetsFS))
	mux.Handle("GET /png", secureHeadersMiddleware(loggingMiddleware(router.downloadPNGHandler)))
	mux.Handle("GET /pdf", secureHeadersMiddleware(loggingMiddleware(router.downloadPDFHandler)))
	mux.Handle("GET /privacy", secureHeadersMiddleware(loggingMiddleware(router.privacyHandler)))
	mux.Handle("GET /attributions", secureHeadersMiddleware(loggingMiddleware(router.attributionsHandler)))
	mux.Handle("GET /", secureHeadersMiddleware(loggingMiddleware(router.mainHandler)))

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
