package session_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/acoshift/session"
)

func mockHandlerFunc(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r.Context())
	s.Set("test", 1)
	w.Write([]byte("ok"))
}

var mockHandler = http.HandlerFunc(mockHandlerFunc)

func TestPanicConfig(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("expected panic when misconfig; but not")
		}
	}()
	session.Middleware(session.Config{})
}

func TestDefautConfig(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{},
	})(mockHandler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	cookie := w.Header().Get("Set-Cookie")
	if len(cookie) == 0 {
		t.Fatalf("expected cookie not empty; got empty")
	}
}

func TestEmptySession(t *testing.T) {
	h := session.Middleware(session.Config{
		Store: &mockStore{
			GetFunc: func(key string) ([]byte, error) {
				t.Fatalf("expected get was not called")
				return nil, nil
			},
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				t.Fatalf("expected set was not called")
				return nil
			},
			DelFunc: func(key string) error {
				t.Fatalf("expected del was not called")
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	cookie := w.Header().Get("Set-Cookie")
	if len(cookie) == 0 {
		t.Fatalf("expected cookie not empty; got empty")
	}
}

func TestSessionSetInStore(t *testing.T) {
	var (
		setCalled bool
		setKey    string
		setValue  []byte
		setTTL    time.Duration
	)

	h := session.Middleware(session.Config{
		MaxAge: time.Second,
		Store: &mockStore{
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				setCalled = true
				setKey = key
				setValue = value
				setTTL = ttl
				return nil
			},
		},
	})(mockHandler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	if !setCalled {
		t.Fatalf("expected store was called; but not")
	}
	if len(setKey) == 0 {
		t.Fatalf("expected key not empty; got empty")
	}
	if len(setValue) == 0 {
		t.Fatalf("expected value not empty; got empty")
	}
	if setTTL != time.Second {
		t.Fatalf("expected ttl to be 1s; got %v", setTTL)
	}

	cs := w.Result().Cookies()
	if len(cs) != 1 {
		t.Fatalf("expected response has 1 cookie; got %d", len(cs))
	}
	if cs[0].Value == setKey {
		t.Fatalf("expected session id was hashed")
	}
}

func TestSessionGetSet(t *testing.T) {
	var (
		setCalled int
		setKey    string
		setValue  []byte
	)

	h := session.Middleware(session.Config{
		MaxAge:       time.Second,
		DisableRenew: true,
		Store: &mockStore{
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				setCalled++
				setKey = key
				setValue = value
				return nil
			},
			GetFunc: func(key string) ([]byte, error) {
				if key != setKey {
					t.Fatalf("expected get key \"%s\"; got \"%s\"", setKey, key)
				}
				return setValue, nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		c, _ := s.Get("test").(int)
		s.Set("test", c+1)
		fmt.Fprintf(w, "%d", c)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	if w.Body.String() != "0" {
		t.Fatalf("expected response to be 0; got %s", w.Body.String())
	}

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Body.String() != "1" {
		t.Fatalf("expected response to be 1; got %s", w.Body.String())
	}

	if setCalled != 2 {
		t.Fatalf("expected store set 2 times; but got %d times", setCalled)
	}
}

func TestSecureFlag(t *testing.T) {
	cases := []struct {
		tls      bool
		flag     session.Secure
		expected bool
	}{
		{false, session.NoSecure, false},
		{false, session.ForceSecure, true},
		{false, session.PreferSecure, false},
		{true, session.NoSecure, false},
		{true, session.ForceSecure, true},
		{true, session.PreferSecure, true},
	}

	for _, c := range cases {
		h := session.Middleware(session.Config{
			Store:  &mockStore{},
			Secure: c.flag,
		})(mockHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		if c.tls {
			r.Header.Set("X-Forwarded-Proto", "https")
		}
		h.ServeHTTP(w, r)

		cs := w.Result().Cookies()
		if len(cs) != 1 {
			t.Fatalf("expected response has 1 cookie; got %d", len(cs))
		}
		if cs[0].Secure != c.expected {
			srv := "http"
			if c.tls {
				srv += "s"
			}
			t.Fatalf("expected cookie secure flag %d for %s to be %v; got %v", c.flag, srv, c.expected, cs[0].Secure)
		}
	}
}

func TestHttpOnlyFlag(t *testing.T) {
	cases := []struct {
		flag bool
	}{
		{false},
		{true},
	}

	for _, c := range cases {
		h := session.Middleware(session.Config{
			Store:    &mockStore{},
			HTTPOnly: c.flag,
		})(mockHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		h.ServeHTTP(w, r)

		cs := w.Result().Cookies()
		if len(cs) != 1 {
			t.Fatalf("expected response has 1 cookie; got %d", len(cs))
		}
		if cs[0].HttpOnly != c.flag {
			t.Fatalf("expected HttpOnly flag to be %v; got %v", c.flag, cs[0].HttpOnly)
		}
	}
}

func TestRotate(t *testing.T) {
	c := 0

	var (
		setCalled int
		setKey    string
		setValue  []byte
	)

	h := session.Middleware(session.Config{
		Store: &mockStore{
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				setCalled++
				if c == 0 {
					setKey = key
					setValue = value
					return nil
				}
				if key == setKey {
					t.Fatalf("expected key after rotate to renew")
				}
				return nil
			},
			GetFunc: func(key string) ([]byte, error) {
				return setValue, nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else {
			s.Rotate()
		}
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if setCalled != 3 {
		t.Fatalf("expected set was called 3 times; got %d times", setCalled)
	}
}

func TestDestroy(t *testing.T) {
	c := 0

	var (
		delCalled bool
		setKey    string
		setValue  []byte
	)

	h := session.Middleware(session.Config{
		Store: &mockStore{
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				setKey = key
				setValue = value
				return nil
			},
			GetFunc: func(key string) ([]byte, error) {
				return setValue, nil
			},
			DelFunc: func(key string) error {
				delCalled = true
				if key != setKey {
					t.Fatalf("expected destroy old key")
				}
				return nil
			},
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		if c == 0 {
			s.Set("test", 1)
			c = 1
		} else {
			s.Destroy()
		}
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if !delCalled {
		t.Fatalf("expected del was called")
	}
}

func TestDisableHashID(t *testing.T) {
	var setKey string

	h := session.Middleware(session.Config{
		DisableHashID: true,
		Store: &mockStore{
			SetFunc: func(key string, value []byte, ttl time.Duration) error {
				setKey = key
				return nil
			},
		},
	})(mockHandler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	cs := w.Result().Cookies()
	if len(cs) != 1 {
		t.Fatalf("expected response has 1 cookie; got %d", len(cs))
	}
	if cs[0].Value != setKey {
		t.Fatalf("expected session id was not hashed")
	}
}
