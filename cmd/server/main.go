package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robyparr/barmycodes/internal/web"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := web.NewRouter(time.Now)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
