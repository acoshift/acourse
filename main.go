package main

import (
	"embed"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"

	_ "github.com/lib/pq"

	"github.com/acoshift/acourse/internal/app"
	"github.com/acoshift/acourse/internal/pkg/auth"
	"github.com/acoshift/acourse/internal/pkg/config"
	"github.com/acoshift/acourse/internal/pkg/email"
	"github.com/acoshift/acourse/internal/pkg/file"
	"github.com/acoshift/acourse/internal/pkg/notify"
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
	config.Init()
	defer config.Close()

	auth.Init()
	email.Init()
	file.Init()
	notify.Init()

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
