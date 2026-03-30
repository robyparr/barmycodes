// +heroku goVersion 1.21
// +heroku install ./cmd/server

module github.com/robyparr/barmycodes

go 1.25.0

require (
	github.com/boombuler/barcode v1.0.2
	github.com/go-pdf/fpdf v0.9.0
	golang.org/x/image v0.38.0
)

require golang.org/x/text v0.35.0 // indirect
