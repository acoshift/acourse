package ds

import (
	"context"

	"cloud.google.com/go/datastore"
)

// PutModel puts a model to datastore
func (client *Client) PutModel(ctx context.Context, src interface{}) error {
	key := src.(KeyGetSetter).GetKey()
	key, err := client.Put(ctx, key, src)
	SetKey(key, src)
	if client.Cache != nil {
		client.Cache.Del(key)
	}
	if err != nil {
		return err
	}
	return nil
}

// PutModels puts models to datastore
func (client *Client) PutModels(ctx context.Context, src interface{}) error {
	xs := valueOf(src)
	keys := make([]*datastore.Key, xs.Len())
	for i := range keys {
		x := xs.Index(i).Interface()
		keys[i] = x.(KeyGetter).GetKey()
	}

	var err error
	l := len(keys)
	p := 500
	if l > p {
		ks := make([]*datastore.Key, 0, l)
		for i := 0; i < l/p+1; i++ {
			m := (i + 1) * p
			if m > l {
				m = l
			}
			if i*p == m {
				break
			}
			k, e := client.PutMulti(ctx, keys[i*p:m], xs.Slice(i*p, m).Interface())
			ks = append(ks, k...)
			if e != nil {
				if err == nil {
					err = e
				} else {
					if errs, ok := err.(datastore.MultiError); ok {
						err = append(errs, e)
					} else {
						err = datastore.MultiError{err, e}
					}
				}
			}
		}
		keys = ks
	} else {
		keys, err = client.PutMulti(ctx, keys, src)
	}
	SetKeys(keys, src)
	if client.Cache != nil {
		client.Cache.DelMulti(keys)
	}
	if err != nil {
		return err
	}
	return nil
}
