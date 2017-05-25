package ds

import (
	"testing"
)

func TestSaveModel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping save model")
	}
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	x := &ExampleModel{Name: "Test1", Value: 1}
	err = client.SaveModel(ctx, x)
	if err != nil {
		t.Fatal(err)
	}
	if x.GetKey() == nil {
		t.Fatalf("expetect key to be assigned")
	}
	if x.CreatedAt.IsZero() || x.UpdatedAt.IsZero() {
		t.Fatalf("expetect stamp model to be assigned")
	}
	if x.ID() == 0 {
		t.Fatalf("expected id to be assigned")
	}
	err = client.DeleteModel(ctx, x)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveModels(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping save models")
	}
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	xs := []*ExampleModel{
		{Name: "Test1", Value: 1},
		{Name: "Test2", Value: 2},
	}
	err = client.SaveModels(ctx, xs)
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range xs {
		if x.GetKey() == nil {
			t.Fatalf("expetect key to be assigned")
		}
		if x.CreatedAt.IsZero() || x.UpdatedAt.IsZero() {
			t.Fatalf("expetect stamp model to be assigned")
		}
	}
	err = client.DeleteModels(ctx, xs)
	if err != nil {
		t.Fatal(err)
	}
}
