package session_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/acoshift/session"
)

func TestManagerGetSave(t *testing.T) {
	var (
		setKey   string
		setValue []byte
	)

	m := session.New(session.Config{
		MaxAge:       time.Second,
		DisableRenew: true,
		Store: &mockStore{
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				setKey = key
				setValue = value
				return nil
			},
			GetFunc: func(key string) ([]byte, error) {
				assert.Equal(t, setKey, key)
				return setValue, nil
			},
		},
	})

	h := func(w http.ResponseWriter, r *http.Request) {
		s := m.Get(r, sessName)
		assert.NotEmpty(t, s.ID())
		c, _ := s.Get("test").(int)
		s.Set("test", c+1)
		fmt.Fprintf(w, "%d", c)

		m.Save(w, s)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h(w, r)

	assert.Equal(t, "0", w.Body.String())

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	h(w, r)
	assert.Equal(t, "1", w.Body.String())
}
