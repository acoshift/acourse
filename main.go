package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"

	_ "github.com/lib/pq"

	"github.com/acoshift/acourse/internal/app"
	"github.com/acoshift/acourse/internal/pkg/config"
)

func main() {
	defer config.Close()

	app := app.New()

	err := app.
		Address(":8080").
		ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
