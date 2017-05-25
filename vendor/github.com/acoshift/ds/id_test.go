package ds

import (
	"testing"
)

func TestParseID(t *testing.T) {
	cases := []struct {
		in  string
		out int64
	}{
		{"0", 0},
		{"-111", -111},
		{"111", 111},
		{"abc", 0},
	}

	for _, c := range cases {
		out := parseID(c.in)
		if out != c.out {
			t.Fatalf("expected parseID %s to be %d; got %d", c.in, c.out, out)
		}
	}
}

func TestBuildIDKeys(t *testing.T) {
	ids := []int64{1, 2, 3, 4, 5}
	keys := BuildIDKeys(kind, ids)
	for i, key := range keys {
		if key.ID != ids[i] {
			t.Fatalf("expected key id to be %d; got %d", ids[i], key.ID)
		}
		if key.Kind != kind {
			t.Fatalf("expected key kind to be %s; got %s", kind, key.Kind)
		}
	}
}

func TestBuildStringIDKeys(t *testing.T) {
	ids := []string{"aa", "bb", "cccc", "123", "456"}
	out := []int64{123, 456}
	keys := BuildStringIDKeys(kind, ids)
	if len(keys) != len(out) {
		t.Fatalf("expected result keys length to be %d; got %d", len(out), len(keys))
	}
	for i, key := range keys {
		if key.ID != out[i] {
			t.Fatalf("expected key id to be %d; got %d", out[i], key.ID)
		}
		if key.Kind != kind {
			t.Fatalf("expected key kind to be %s; got %s", kind, key.Kind)
		}
	}
}

func TestBuildNameIDKeys(t *testing.T) {
	ids := []string{"aa", "bb", "cccc", "123", "456"}
	keys := BuildNameKeys(kind, ids)
	if len(keys) != len(ids) {
		t.Fatalf("expected result keys length to be %d; got %d", len(ids), len(keys))
	}
	for i, key := range keys {
		if key.Name != ids[i] {
			t.Fatalf("expected key id to be %s; got %s", ids[i], key.Name)
		}
		if key.Kind != kind {
			t.Fatalf("expected key kind to be %s; got %s", kind, key.Kind)
		}
	}
}

func TestExtractKey(t *testing.T) {
	x := &ExampleModel{}
	x.SetID(kind, 10)
	key := ExtractKey(x)
	if key == nil {
		t.Fatalf("expected key not nil")
	}
}

func TestExtractKeys(t *testing.T) {
	xs := make([]*ExampleModel, 10)
	for i := range xs {
		xs[i] = &ExampleModel{}
		xs[i].SetID(kind, int64(i))
	}
	keys := ExtractKeys(xs)
	for i, key := range keys {
		if key == nil {
			t.Fatalf("expected key not nil")
		}
		if key.ID != int64(i) {
			t.Fatalf("expected key id to be %d; got %d", i, key.ID)
		}
	}
}

func TestExtractKeysPtr(t *testing.T) {
	xs := make([]*ExampleModel, 10)
	for i := range xs {
		xs[i] = &ExampleModel{}
		xs[i].SetID(kind, int64(i))
	}
	keys := ExtractKeys(&xs)
	for i, key := range keys {
		if key == nil {
			t.Fatalf("expected key not nil")
		}
		if key.ID != int64(i) {
			t.Fatalf("expected key id to be %d; got %d", i, key.ID)
		}
	}
}
