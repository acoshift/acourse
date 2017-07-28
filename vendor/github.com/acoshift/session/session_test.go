package session

import (
	"testing"
	"time"
)

func TestEncodeEmpty(t *testing.T) {
	s := Session{}
	b := s.encode()
	if b == nil {
		t.Fatalf("expected encode always not return nil; got nil")
	}
	if len(b) != 0 {
		t.Fatalf("expected encode empty session return empty length slice; got %d length", len(b))
	}
}

func TestEncodeUnregisterType(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("expected encode unregister type panic")
		}
	}()
	type a struct{}
	s := Session{}
	s.Set("a", a{})
	s.encode()
}

func TestSessionOperation(t *testing.T) {
	s := Session{}
	if x := s.Get("a"); x != nil {
		t.Fatalf("expected get data from empty session return nil; got %v", x)
	}
	s.Del("a")
	if s.data != nil {
		t.Fatalf("expected data to be nil; got %v", s.data)
	}
	s.Set("a", 1)
	if x, _ := s.Get("a").(int); x != 1 {
		t.Fatalf("expected get return 1; got %d", x)
	}
	s.Del("a")
	if x := s.Get("a"); x != nil {
		t.Fatalf("expected get data after delete to be nil; got %v", x)
	}
}

func TestShouldRenew(t *testing.T) {
	s := Session{}
	s.Set(timestampKey{}, int64(-1))
	if s.shouldRenew() {
		t.Fatalf("expected sec -1 should not renew")
	}

	s.Set(timestampKey{}, int64(0))
	if !s.shouldRenew() {
		t.Fatalf("expected sec 0 should renew")
	}

	now := time.Now().Unix()

	s.MaxAge = 10 * time.Second
	s.Set(timestampKey{}, now-7)
	if !s.shouldRenew() {
		t.Fatalf("expected sec -7 of max-age 10 should renew")
	}

	s.Set(timestampKey{}, now-5)
	if !s.shouldRenew() {
		t.Fatalf("expected sec -5 of max-age 10 should renew")
	}

	s.Set(timestampKey{}, now-3)
	if s.shouldRenew() {
		t.Fatalf("expected sec -3 of max-age 10 should not renew")
	}
}
