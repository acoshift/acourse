package flash_test

import (
	"testing"

	. "github.com/acoshift/flash"
)

func TestNew(t *testing.T) {
	f := New()
	if f == nil {
		t.Errorf("expected New returns valid flash; got nil")
	}
	if f.Count() > 0 {
		t.Errorf("expected New returns empty flash; got %v", f)
	}
	if f.Changed() {
		t.Errorf("expected New returns unchanged flash; got changed")
	}
}

func TestEncodeDecode(t *testing.T) {
	var (
		f, p *Flash
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
	if f.Changed() {
		t.Errorf("expected Encode empty flash still unchange; got changed")
	}
	p, err = Decode([]byte{})
	if err != nil {
		t.Errorf("expected Decode empty bytes not error; got %v", err)
	}
	if p == nil {
		t.Errorf("expected Decode empty bytes returns non-nil flash; got nil")
	}
	if p.Count() > 0 {
		t.Errorf("expected Decode empty bytes returns empty flash; got %v", nil)
	}
	if f.Changed() {
		t.Errorf("expected Decode empty flash returns unchange flash; got changed")
	}

	f = New()
	f.Add("a", "1")
	if !f.Changed() {
		t.Errorf("expected Add data to empty flash must changed; got unchange")
	}
	b, err = f.Encode()
	if err != nil {
		t.Errorf("expected Encode non-zero length flash not error; got %v", err)
	}
	if len(b) == 0 {
		t.Errorf("expected Encode non-zero length flash returns non-zero length bytes; got %v", b)
	}
	if !f.Changed() {
		t.Errorf("expected Encode changed flash still changed; got unchange")
	}
	p, err = Decode(b)
	if err != nil {
		t.Errorf("expected Decode not error; got %v", err)
	}
	if n := f.Count(); n != 1 {
		t.Errorf("expected Decode returns 1 key flash; got %d", n)
	}
	if f.Values()["a"][0] != "1" {
		t.Errorf("expected Decode returns same value; got %v", f.Values()["a"][0])
	}
	if p.Changed() {
		t.Errorf("expected Decode not empty flash must unchange; got changed")
	}
	p.Clear()
	if !p.Changed() {
		t.Errorf("expected clear not empty unchanged flash must changed; got unchange")
	}

	p, _ = Decode(b)
	p.Del("a")
	if !p.Changed() {
		t.Errorf("expected Del from not empty unchanged flash must change; got unchanged")
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
	if p := len(f.Values()["c"]); p != 2 {
		t.Errorf("expected f.c has 2 values; got %d value(s)", p)
	}

	f.Del("b")
	if f.Has("b") {
		t.Errorf("expected f don't have b")
	}

	f.Del("empty")
	if f.Has("empty") {
		t.Errorf("expected f don't have empty")
	}

	f.Clear()
	if f.Count() > 0 {
		t.Errorf("expected f empty after clear")
	}
}

func TestClone(t *testing.T) {
	f := New()
	f.Add("a", "1")
	f.Add("a", "2")
	f.Add("b", "3")

	p := f.Clone()

	fb := f.Values().Encode()
	pb := p.Values().Encode()
	if fb != pb {
		t.Fatalf("expected cloned encode to be same value")
	}

	f.Clear()
	fb = f.Values().Encode()
	pb = p.Values().Encode()
	if fb == pb {
		t.Fatalf("expected cloned encode and cleard original not same value")
	}
}

func TestDecodeError(t *testing.T) {
	f, err := Decode([]byte("invalid data"))
	if err == nil {
		t.Fatalf("expected decode invalid data returns error; got nil")
	}
	if f != nil {
		t.Fatalf("expected decode invalid data returns nil flash; got %v", f)
	}
}
