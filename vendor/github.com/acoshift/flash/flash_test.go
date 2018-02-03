package flash

import (
	"testing"
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

	if f.Has("a") {
		t.Errorf("expected has empty flash return false; got true")
	}
	if f.Get("a") != nil {
		t.Errorf("expected get from empty flash return nil")
	}
	l := f.Value("a")
	if l == nil {
		t.Errorf("expected value always return non-nil slice")
	}
	if len(l) > 0 {
		t.Errorf("expected value from empty flash return zero-length slice")
	}

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

	f.Del("b")
	if f.Has("b") {
		t.Errorf("expected f don't have b")
	}

	f.Del("empty")
	if f.Has("empty") {
		t.Errorf("expected f don't have empty")
	}

	l = f.Value("c")
	if len(l) != 2 {
		t.Errorf("expected value from 'c' return 2 length slice; got %d", len(l))
	}
	if l[0].(string) != "3" && l[1].(string) != "4" {
		t.Errorf("expected value return values from added")
	}

	f.Clear()
	if f.Count() > 0 {
		t.Errorf("expected f empty after clear")
	}

	f.Set("string", "string value")
	if f.GetString("string") != "string value" {
		t.Errorf("expected get string return valid value")
	}
	f.Set("int", 1)
	if f.GetInt("int") != 1 {
		t.Errorf("expected get int return valid value")
	}
	f.Set("int64", int64(1))
	if f.GetInt64("int64") != int64(1) {
		t.Errorf("expected get int64 return valid value")
	}
	f.Set("float32", float32(1.2))
	if f.GetFloat32("float32") != float32(1.2) {
		t.Errorf("expected get float32 return valid value")
	}
	f.Set("float64", float64(1.5))
	if f.GetFloat64("float64") != float64(1.5) {
		t.Errorf("expected get float64 return valid value")
	}
	f.Set("bool", true)
	if f.GetBool("bool") != true {
		t.Errorf("expected get bool return valid value")
	}

	if len(f.Value("empty")) != 0 {
		t.Errorf("expected value from empty key return zero-length slice")
	}

	v := f.Values()
	if len(v) != f.Count() {
		t.Errorf("expected Values()'s length equals to Count()")
	}

	f.v["a"] = []interface{}{}
	if f.Get("a") != nil {
		t.Errorf("expected get from non-nil, zero-length key return nil")
	}
}

func equals(f *Flash, p *Flash) bool {
	if f.Count() != p.Count() {
		return false
	}
	for k := range f.v {
		if len(f.v[k]) != len(p.v[k]) {
			return false
		}
		for kk := range f.v[k] {
			if f.v[k][kk] != p.v[k][kk] {
				return false
			}
		}
	}
	return true
}

func TestClone(t *testing.T) {
	f := New()
	f.Add("a", "1")
	f.Add("a", "2")
	f.Add("b", "3")

	p := f.Clone()

	if f == p {
		t.Errorf("expected cloned flash don't have same pointer")
	}

	if !equals(f, p) {
		t.Errorf("expected cloned encode to be same value")
	}

	f.Clear()
	if equals(f, p) {
		t.Errorf("expected cloned encode and cleard original not same value")
	}
}

func TestEncodeError(t *testing.T) {
	v := struct{}{}
	f := New()
	f.Set("key", &v)
	b, err := f.Encode()
	if err == nil {
		t.Errorf("expected encode unregistered gob struct error; got nil")
	}
	if b == nil {
		t.Errorf("expected result of encode error not nil")
	}
	if len(b) > 0 {
		t.Errorf("expected length of encode error to be 0; got nil")
	}
}

func TestDecodeError(t *testing.T) {
	f, err := Decode([]byte("invalid data"))
	if err == nil {
		t.Errorf("expected decode invalid data returns error; got nil")
	}
	if f != nil {
		t.Errorf("expected decode invalid data returns nil flash; got %v", f)
	}
}

func TestClearEmpty(t *testing.T) {
	f := New()
	f.Clear()
	if f.Changed() {
		t.Errorf("expected clear empty flash must not changed")
	}
}
