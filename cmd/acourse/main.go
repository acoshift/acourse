package main

import (
	"log"
	"net/http"

	"github.com/acoshift/acourse/pkg/app"
	_ "github.com/lib/pq"
)

func main() {
	// lets reverse proxy handle other settings
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Handler,
	}

	log.Println("Start server at :8080")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
