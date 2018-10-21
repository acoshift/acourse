package admin

import (
	"net/http"
	"time"

	"github.com/acoshift/methodmux"
	"github.com/moonrhythm/hime"
)

// Config is admin config
type Config struct {
	Location *time.Location
}

// New creates admin handler
func New(cfg Config) http.Handler {
	c := &ctrl{cfg}

	mux := http.NewServeMux()

	mux.Handle("/users", methodmux.Get(
		hime.Handler(c.users),
	))
	mux.Handle("/courses", methodmux.Get(
		hime.Handler(c.courses),
	))
	mux.Handle("/payments/pending", methodmux.GetPost(
		hime.Handler(c.pendingPayments),
		hime.Handler(c.postPendingPayment),
	))
	mux.Handle("/payments/history", methodmux.Get(
		hime.Handler(c.historyPayments),
	))
	mux.Handle("/payments/reject", methodmux.GetPost(
		hime.Handler(c.rejectPayment),
		hime.Handler(c.postRejectPayment),
	))

	return mux
}

type ctrl struct {
	Config
}
