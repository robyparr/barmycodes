package web

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/robyparr/barmycodes/internal"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed assets/*
var assetsFS embed.FS
var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.Handle("GET /assets/", http.FileServerFS(assetsFS))
	router.HandleFunc("GET /", mainHandler)

	return router
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	barcodeType := r.URL.Query().Get("type")
	barcodeValues := r.URL.Query()["b[]"]
	barcodes, _ := internal.GenerateBarcodes(barcodeValues, barcodeType)

	tmpl.ExecuteTemplate(w, "index.html.tmpl", barcodes)
}
