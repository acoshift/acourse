package flash_test

import (
	"bytes"
	"testing"

	. "github.com/acoshift/flash"
)

func TestNew(t *testing.T) {
	f := New()
	if f == nil {
		t.Errorf("expected New returns valid flash; got nil")
	}
	if len(f) > 0 {
		t.Errorf("expected New returns empty flash; got %v", f)
	}
}

func TestEncodeDecode(t *testing.T) {
	var (
		f, p Flash
		b    []byte
		err  error
	)

	f = New()
	b, err = f.Encode()
	if err != nil {
		t.Errorf("expected Encode empty flash not error; got %v", err)
	}
	if len(b) > 0 {
		t.Errorf("expected Encode empty flash returns zero length bytes; got %v", b)
	}
	p, err = Decode([]byte{})
	if err != nil {
		t.Errorf("expected Decode empty bytes not error; got %v", err)
	}
	if p == nil {
		t.Errorf("expected Decode empty bytes returns non-nil flash; got nil")
	}
	if len(p) > 0 {
		t.Errorf("expected Decode empty bytes returns empty flash; got %v", nil)
	}

	f = New()
	f.Add("a", "1")
	b, err = f.Encode()
	if err != nil {
		t.Errorf("expected Encode non-zero length flash not error; got %v", err)
	}
	if len(b) == 0 {
		t.Errorf("expected Encode non-zero length flash returns non-zero length bytes; got %v", b)
	}
	p, err = Decode(b)
	if err != nil {
		t.Errorf("expected Decode not error; got %v", err)
	}
	if n := len(f); n != 1 {
		t.Errorf("expected Decode returns 1 key flash; got %d", n)
	}
	if f["a"][0] != "1" {
		t.Errorf("expected Decode returns same value; got %v", f["a"][0])
	}
}

func TestOperators(t *testing.T) {
	f := New()
	f.Set("a", "1")
	f.Set("b", "2")
	f.Add("c", "3")
	f.Add("c", "4")

	if !f.Has("a") {
		t.Errorf("expected f has 'a'")
	}
	if !f.Has("c") {
		t.Errorf("expected f has 'c'")
	}
	if p := f.Get("a"); p != "1" {
		t.Errorf("expected Get 'a' from f is 1; got %s", p)
	}
	if p := f.Get("c"); p != "3" {
		t.Errorf("expected Get 'c' from f is 3; got %s", p)
	}
	if p := len(f["c"]); p != 2 {
		t.Errorf("expected f.c has 2 values; got %d value(s)", p)
	}

	f.Del("b")
	if f.Has("b") {
		t.Errorf("expected f don't have b")
	}

	f.Clear()
	if len(f) > 0 {
		t.Errorf("expected f empty after clear")
	}
}

func TestClone(t *testing.T) {
	f := New()
	f.Add("a", "1")
	f.Add("a", "2")
	f.Add("b", "3")

	p := f.Clone()

	fb, _ := f.Encode()
	pb, _ := p.Encode()
	if !bytes.Equal(fb, pb) {
		t.Fatalf("expected cloned encode to be same value")
	}

	f.Clear()
	fb, _ = f.Encode()
	pb, _ = p.Encode()
	if bytes.Equal(fb, pb) {
		t.Fatalf("expected cloned encode and cleard original not same value")
	}
}
