package app

import (
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/admin"
	"github.com/acoshift/acourse/internal/app/app"
	"github.com/acoshift/acourse/internal/app/auth"
	"github.com/acoshift/acourse/internal/app/editor"
	"github.com/acoshift/acourse/internal/app/view"
)

// Handler creates new app's handler
func Handler(baseURL string, loc *time.Location) http.Handler {
	methodmux.FallbackHandler = hime.Handler(view.NotFound)

	m := http.NewServeMux()
	app.Mount(m, baseURL)
	auth.Mount(m)
	editor.Mount(m)
	admin.Mount(m, loc)

	return m
}
