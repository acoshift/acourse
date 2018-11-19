package admin

import (
	"net/http"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/pkg/context/appctx"
)

// Mount mounts admin handlers
func Mount(m *http.ServeMux) {
	mux := http.NewServeMux()
	mux.Handle("/admin/users", methodmux.Get(
		hime.Handler(getUsers),
	))
	mux.Handle("/admin/courses", methodmux.Get(
		hime.Handler(getCourses),
	))
	mux.Handle("/admin/payments/pending", methodmux.GetPost(
		hime.Handler(getPendingPayments),
		hime.Handler(postPendingPayment),
	))
	mux.Handle("/admin/payments/history", methodmux.Get(
		hime.Handler(getHistoryPayments),
	))
	mux.Handle("/admin/payments/reject", methodmux.GetPost(
		hime.Handler(getRejectPayment),
		hime.Handler(postRejectPayment),
	))

	m.Handle("/admin/", onlyAdmin(mux))
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
