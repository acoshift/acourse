package main

import (
	"embed"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"

	_ "github.com/lib/pq"

	"github.com/acoshift/acourse/internal/app"
	"github.com/acoshift/acourse/internal/pkg/config"
)

//go:embed template/* template/**/*
var templateFS embed.FS

//go:embed settings/routes.yaml
var routes []byte

//go:embed settings/template.yaml
var templateConfig []byte

//go:embed settings/static.yaml
var static []byte

//go:embed assets
var assetsFS embed.FS

func main() {
	defer config.Close()

	assets, _ := fs.Sub(assetsFS, "assets")

	err := app.New(app.Config{
		Routes:         routes,
		Static:         static,
		TemplateConfig: templateConfig,
		Template:       templateFS,
		Assets:         assets,
	}).
		Address(":8080").
		ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
