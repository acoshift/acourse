package memory

import (
	"testing"
	"time"
)

func TestMemory(t *testing.T) {
	s := New(Config{CleanupInterval: 10 * time.Millisecond})
	err := s.Set("a", []byte("test"), time.Millisecond)
	if err != nil {
		t.Fatalf("expected set not error; got %v", err)
	}

	time.Sleep(5 * time.Millisecond)
	b, err := s.Get("a")
	if b != nil {
		t.Fatalf("expected expired key return nil value; got %v", b)
	}
	if err == nil {
		t.Fatalf("expected expired key return error; got nil")
	}

	s.Set("a", []byte("test"), time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, err = s.Get("a")
	if err == nil {
		t.Fatalf("expected expired key return error; got nil")
	}

	s.Set("a", []byte("test"), time.Second)
	b, err = s.Get("a")
	if err != nil {
		t.Fatalf("expected get valid key return not nil error; got %v", err)
	}
	if string(b) != "test" {
		t.Fatalf("expected get return same value as set")
	}

	s.Del("a")
	_, err = s.Get("a")
	if err == nil {
		t.Fatalf("expected get deleted key to return error; got nil")
	}
}
