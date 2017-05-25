package ds

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
)

func parseID(id string) int64 {
	r, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}
	return r
}

// BuildIDKeys builds datastore keys from id keys
func BuildIDKeys(kind string, ids []int64) []*datastore.Key {
	keys := make([]*datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = datastore.IDKey(kind, id, nil)
	}
	return keys
}

// BuildStringIDKeys builds datastore keys from string id keys
func BuildStringIDKeys(kind string, ids []string) []*datastore.Key {
	keys := make([]*datastore.Key, 0, len(ids))
	for _, id := range ids {
		if tid := parseID(id); tid != 0 {
			keys = append(keys, datastore.IDKey(kind, tid, nil))
		}
	}
	return keys
}

// BuildNameKeys builds datastore keys from name keys
func BuildNameKeys(kind string, names []string) []*datastore.Key {
	keys := make([]*datastore.Key, len(names))
	for i, name := range names {
		keys[i] = datastore.NameKey(kind, name, nil)
	}
	return keys
}

// ExtractKey returns key from model
func ExtractKey(src interface{}) *datastore.Key {
	return src.(KeyGetter).GetKey()
}

// ExtractKeys returns keys from models
func ExtractKeys(src interface{}) []*datastore.Key {
	xs := valueOf(src)
	l := xs.Len()
	keys := make([]*datastore.Key, l)
	for i := 0; i < l; i++ {
		keys[i] = ExtractKey(xs.Index(i).Interface())
	}
	return keys
}

// AllocateIDModel allocates id for model
func (client *Client) AllocateIDModel(ctx context.Context, kind string, src interface{}) error {
	m := src.(KeyGetSetter)
	if m.GetKey() == nil {
		m.SetKey(datastore.IncompleteKey(kind, nil))
	}
	keys, err := client.AllocateIDs(ctx, []*datastore.Key{m.GetKey()})
	if err != nil {
		return err
	}
	m.SetKey(keys[0])
	return nil
}

// AllocateIDModels allocates id for models
func (client *Client) AllocateIDModels(ctx context.Context, kind string, src interface{}) error {
	xs := valueOf(src)
	keys := make([]*datastore.Key, xs.Len())
	for i := range keys {
		x := xs.Index(i).Interface()
		m := x.(KeyGetSetter)
		if m.GetKey() == nil {
			m.SetKey(datastore.IncompleteKey(kind, nil))
		}
		keys[i] = x.(KeyGetter).GetKey()
	}
	keys, err := client.AllocateIDs(ctx, keys)
	if err != nil {
		return err
	}
	SetKeys(keys, src)
	return nil
}
