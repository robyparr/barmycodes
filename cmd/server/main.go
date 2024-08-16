package main

import (
	"log"
	"net/http"

	"github.com/robyparr/barmycodes/internal/web"
)

func main() {
	address := "127.0.0.1:8080"
	router := web.NewRouter()
	log.Fatal(http.ListenAndServe(address, router))
}
