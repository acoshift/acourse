package admin

import (
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

// Mount mounts admin handlers
func Mount(m *http.ServeMux, loc *time.Location) {
	c := &ctrl{loc}

	mux := http.NewServeMux()
	mux.Handle("/admin/users", methodmux.Get(
		hime.Handler(c.getUsers),
	))
	mux.Handle("/admin/courses", methodmux.Get(
		hime.Handler(c.getCourses),
	))
	mux.Handle("/admin/payments/pending", methodmux.GetPost(
		hime.Handler(c.getPendingPayments),
		hime.Handler(c.postPendingPayment),
	))
	mux.Handle("/admin/payments/history", methodmux.Get(
		hime.Handler(c.getHistoryPayments),
	))
	mux.Handle("/admin/payments/reject", methodmux.GetPost(
		hime.Handler(c.getRejectPayment),
		hime.Handler(c.postRejectPayment),
	))

	m.Handle("/admin/", onlyAdmin(mux))
}

type ctrl struct {
	Location *time.Location
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
