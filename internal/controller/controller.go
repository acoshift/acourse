package controller

import (
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/context/appctx"
	"github.com/acoshift/acourse/internal/controller/admin"
	"github.com/acoshift/acourse/internal/controller/app"
	"github.com/acoshift/acourse/internal/controller/auth"
	"github.com/acoshift/acourse/internal/controller/editor"
	"github.com/acoshift/acourse/internal/controller/share"
)

// Mount mounts controllers into mux
func Mount(m *http.ServeMux, baseURL string, loc *time.Location) {
	methodmux.FallbackHandler = hime.Handler(share.NotFound)

	m.Handle("/", app.New(app.Config{BaseURL: baseURL}))
	m.Handle("/auth/", http.StripPrefix("/auth", notSignedIn(auth.New())))
	m.Handle("/editor/", http.StripPrefix("/editor", editor.New()))
	m.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin.New(admin.Config{Location: loc}))))
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
