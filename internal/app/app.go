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
	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

// Handler creates new app's handler
func Handler(baseURL string, loc *time.Location) http.Handler {
	methodmux.FallbackHandler = hime.Handler(view.NotFound)

	m := http.NewServeMux()
	m.Handle("/", app.New(app.Config{BaseURL: baseURL}))
	m.Handle("/auth/", http.StripPrefix("/auth", notSignedIn(auth.New())))
	m.Handle("/editor/", http.StripPrefix("/editor", editor.New()))
	m.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(admin.New(admin.Config{Location: loc}))))

	return m
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
