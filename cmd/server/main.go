package main

import (
	"log"
	"net/http"
	"time"

	"github.com/robyparr/barmycodes/internal/web"
)

func main() {
	address := "127.0.0.1:8080"
	router := web.NewRouter(time.Now)
	log.Fatal(http.ListenAndServe(address, router))
}
