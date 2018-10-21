package controller

import (
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/acoshift/middleware"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/controller/admin"
	"github.com/acoshift/acourse/controller/app"
	"github.com/acoshift/acourse/controller/auth"
	"github.com/acoshift/acourse/controller/editor"
	"github.com/acoshift/acourse/controller/share"
)

// Mount mounts controllers into mux
func Mount(m *http.ServeMux, baseURL string, loc *time.Location) {
	methodmux.FallbackHandler = hime.Handler(share.NotFound)

	m.Handle("/", app.New(app.Config{
		BaseURL: baseURL,
	}))

	m.Handle("/auth/", http.StripPrefix("/auth", middleware.Chain(
		notSignedIn,
	)(auth.New())))

	m.Handle("/editor/", http.StripPrefix("/editor", middleware.Chain(
	// -
	)(editor.New())))

	m.Handle("/admin/", http.StripPrefix("/admin", middleware.Chain(
		onlyAdmin,
	)(admin.New(admin.Config{
		Location: loc,
	}))))
}

func onlyAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := appctx.GetUser(r.Context())
		if u == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.Role.Admin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func notSignedIn(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := appctx.GetUserID(r.Context())
		if len(id) > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
