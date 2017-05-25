package ds

import (
	"testing"
)

func TestGetByKey(t *testing.T) {
	skipShort(t, "GetByKey")
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}

	keys := prepareData(client)
	defer removeData(client)

	var x ExampleModel

	err = client.GetByKey(ctx, keys[0], &x)
	if err != nil {
		t.Fatal(err)
	}
	if !x.GetKey().Equal(keys[0]) {
		t.Fatalf("key not equals")
	}

	xs := make([]*ExampleModel, len(keys))
	err = client.GetByKeys(ctx, keys, xs)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != len(xs) {
		t.Fatalf("keys and result len not equals")
	}
	for i := range keys {
		if !keys[i].Equal(xs[i].GetKey()) {
			t.Fatalf("key not equals")
		}
	}

	var xs2 []*ExampleModel
	err = client.GetByKeys(ctx, keys, &xs2)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != len(xs2) {
		t.Fatalf("keys and result len not equals")
	}
	for i := range keys {
		if !keys[i].Equal(xs2[i].GetKey()) {
			t.Fatalf("key not equals")
		}
	}
}

func TestGetByModel(t *testing.T) {
	skipShort(t, "GetByModel")
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}

	keys := prepareData(client)
	defer removeData(client)
	var x ExampleModel
	x.Key = keys[0]

	err = client.GetByModel(ctx, &x)
	if err != nil {
		t.Fatal(err)
	}
	if !x.GetKey().Equal(keys[0]) {
		t.Fatalf("key not equals")
	}

	xs := make([]*ExampleModel, len(keys))
	for i, key := range keys {
		xs[i] = &ExampleModel{}
		xs[i].Key = key
	}
	err = client.GetByModels(ctx, xs)
	if err != nil {
		t.Fatal(err)
	}
	for i := range keys {
		if !keys[i].Equal(xs[i].GetKey()) {
			t.Fatalf("key not equals")
		}
	}

	xs2 := make([]*ExampleModel, len(keys))
	for i, key := range keys {
		xs2[i] = &ExampleModel{}
		xs2[i].Key = key
	}
	err = client.GetByModels(ctx, &xs2)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != len(xs2) {
		t.Fatalf("keys and result len not equals")
	}
	for i := range keys {
		if !keys[i].Equal(xs2[i].GetKey()) {
			t.Fatalf("key not equals")
		}
	}
}
