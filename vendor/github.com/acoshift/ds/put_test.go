package ds

import (
	"testing"

	"cloud.google.com/go/datastore"
)

func TestPutModel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping put model")
	}
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	x := &ExampleModel{Name: "Test1", Value: 1}
	err = client.PutModel(ctx, x)
	if err != datastore.ErrInvalidKey {
		t.Fatalf("expected error to be %v; got %v", datastore.ErrInvalidKey, err)
	}
	x.SetID(kind, 99)
	err = client.PutModel(ctx, x)
	if err != nil {
		t.Fatal(err)
	}
	if !x.CreatedAt.IsZero() || !x.UpdatedAt.IsZero() {
		t.Fatalf("expetect stamp model not assigned")
	}
	err = client.DeleteModel(ctx, x)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPutModels(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping put model")
	}
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	xs := []*ExampleModel{
		&ExampleModel{Name: "Test1", Value: 1},
		&ExampleModel{Name: "Test2", Value: 2},
	}
	err = client.PutModels(ctx, xs)
	if err == nil {
		t.Fatalf("expected error not nil")
	}
	for i, x := range xs {
		x.SetID(kind, int64(i+100))
	}
	err = client.PutModels(ctx, xs)
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range xs {
		if !x.CreatedAt.IsZero() || !x.UpdatedAt.IsZero() {
			t.Fatalf("expetect stamp model not assigned")
		}
	}
	err = client.DeleteModels(ctx, xs)
	if err != nil {
		t.Fatal(err)
	}
}
