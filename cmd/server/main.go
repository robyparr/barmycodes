package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/robyparr/barmycodes/internal"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed assets/*
var assetsFS embed.FS
var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))

func handler(w http.ResponseWriter, r *http.Request) {
	barcodeType := r.URL.Query().Get("type")
	barcodeValues := r.URL.Query()["b[]"]
	barcodes, _ := internal.GenerateBarcodes(barcodeValues, barcodeType)

	tmpl.ExecuteTemplate(w, "index.html.tmpl", barcodes)
}

func main() {
	http.Handle("GET /assets/", http.FileServerFS(assetsFS))
	http.HandleFunc("GET /", handler)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
