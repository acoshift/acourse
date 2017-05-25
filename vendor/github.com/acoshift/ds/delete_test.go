package ds

import (
	"strconv"
	"testing"
)

func TestDeleteByID(t *testing.T) {
	skipShort(t, "DeleteByID")
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	keys := prepareData(client)
	defer removeData(client)
	err = client.DeleteByID(ctx, keys[0].Kind, keys[0].ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteByStringID(t *testing.T) {
	skipShort(t, "DeleteByStringID")
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	keys := prepareData(client)
	defer removeData(client)
	err = client.DeleteByStringID(ctx, keys[0].Kind, strconv.FormatInt(keys[0].ID, 10))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteByIDs(t *testing.T) {
	skipShort(t, "DeleteByIDs")
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	keys := prepareData(client)
	defer removeData(client)
	ids := make([]int64, len(keys))
	for i := range keys {
		ids[i] = keys[i].ID
	}
	err = client.DeleteByIDs(ctx, keys[0].Kind, ids)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteByStringIDs(t *testing.T) {
	skipShort(t, "DeleteByStringIDs")
	client, err := initClient()
	if err != nil {
		t.Fatal(err)
	}
	keys := prepareData(client)
	defer removeData(client)
	ids := make([]string, len(keys))
	for i := range keys {
		ids[i] = strconv.FormatInt(keys[i].ID, 10)
	}
	err = client.DeleteByStringIDs(ctx, keys[0].Kind, ids)
	if err != nil {
		t.Fatal(err)
	}
}
