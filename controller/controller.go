package controller

import (
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/controller/admin"
	"github.com/acoshift/acourse/controller/app"
	"github.com/acoshift/acourse/controller/auth"
	"github.com/acoshift/acourse/controller/editor"
	"github.com/acoshift/acourse/controller/share"
	"github.com/acoshift/acourse/internal"
	"github.com/acoshift/acourse/repository"
)

// Mount mounts controllers into mux
func Mount(m *http.ServeMux, baseURL string, loc *time.Location) {
	methodmux.FallbackHandler = hime.Handler(share.NotFound)

	m.Handle("/", app.New(app.Config{
		BaseURL:    baseURL,
		Repository: repository.NewApp(),
	}))

	m.Handle("/auth/", http.StripPrefix("/auth", middleware.Chain(
		internal.NotSignedIn,
	)(auth.New())))

	m.Handle("/editor/", http.StripPrefix("/editor", middleware.Chain(
	// -
	)(editor.New())))

	m.Handle("/admin/", http.StripPrefix("/admin", middleware.Chain(
		internal.OnlyAdmin,
	)(admin.New(admin.Config{
		Location:   loc,
		Repository: repository.NewAdmin(),
	}))))
}
