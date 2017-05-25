package ds

import (
	"context"
	"reflect"

	"cloud.google.com/go/datastore"
)

// GetByKey retrieves model from datastore by key
func (client *Client) GetByKey(ctx context.Context, key *datastore.Key, dst interface{}) error {
	if client.Cache != nil && client.Cache.Get(key, dst) == nil {
		return nil
	}
	err := client.Get(ctx, key, dst)
	SetKey(key, dst)
	if client.Cache != nil {
		client.Cache.Set(key, dst)
	}
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) getByKeys(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	var err error
	l := len(keys)
	p := 1000
	if l > p {
		rfDst := valueOf(dst)
		for i := 0; i < l/p+1; i++ {
			m := (i + 1) * p
			if m > l {
				m = l
			}
			if i*p == m {
				break
			}
			e := client.GetMulti(ctx, keys[i*p:m], rfDst.Slice(i*p, m).Interface())
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
	} else {
		err = client.GetMulti(ctx, keys, dst)
	}
	SetKeys(keys, dst)
	if client.Cache != nil {
		client.Cache.SetMulti(keys, dst)
	}
	if err != nil {
		return err
	}
	return nil
}

// GetByKeys retrieves models from datastore by keys
func (client *Client) GetByKeys(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	// prepare slice if dst is pointer to 0 len slice
	if rf := reflect.ValueOf(dst); rf.Kind() == reflect.Ptr {
		rs := rf.Elem()
		if rs.Kind() == reflect.Slice && rs.Len() == 0 {
			l := len(keys)
			rs.Set(reflect.MakeSlice(rs.Type(), l, l))
		}
		dst = rs.Interface()
	}

	if len(keys) == 0 {
		return nil
	}

	if client.Cache != nil {
		err := client.Cache.GetMulti(keys, dst)
		if err == nil {
			nfKeys := []*datastore.Key{}
			nfMap := []int{}
			rf := valueOf(dst)
			for i := 0; i < rf.Len(); i++ {
				if rf.Index(i).IsNil() {
					nfKeys = append(nfKeys, keys[i])
					nfMap = append(nfMap, i)
				}
			}
			l := len(nfKeys)
			if l > 0 {
				nfDstRf := reflect.MakeSlice(rf.Type(), l, l)
				err = client.getByKeys(ctx, nfKeys, nfDstRf.Interface())
				for i, k := range nfMap {
					rf.Index(k).Set(nfDstRf.Index(i))
				}
			}
			SetKeys(keys, dst)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return client.getByKeys(ctx, keys, dst)
}

// GetByModel retrieves model from datastore by key from model
func (client *Client) GetByModel(ctx context.Context, dst interface{}) error {
	key := ExtractKey(dst)
	return client.GetByKey(ctx, key, dst)
}

// GetByModels retrieves models from datastore by keys from models
func (client *Client) GetByModels(ctx context.Context, dst interface{}) error {
	keys := ExtractKeys(dst)
	return client.GetByKeys(ctx, keys, dst)
}

// GetByID retrieves model from datastore by id
func (client *Client) GetByID(ctx context.Context, kind string, id int64, dst interface{}) error {
	return client.GetByKey(ctx, datastore.IDKey(kind, id, nil), dst)
}

// GetByIDs retrieves models from datastore by ids
func (client *Client) GetByIDs(ctx context.Context, kind string, ids []int64, dst interface{}) error {
	keys := BuildIDKeys(kind, ids)
	return client.GetByKeys(ctx, keys, dst)
}

// GetByStringID retrieves model from datastore by string id
func (client *Client) GetByStringID(ctx context.Context, kind string, id string, dst interface{}) error {
	tid := parseID(id)
	if tid == 0 {
		return datastore.ErrInvalidKey
	}
	return client.GetByKey(ctx, datastore.IDKey(kind, tid, nil), dst)
}

// GetByStringIDs retrieves models from datastore by string ids
func (client *Client) GetByStringIDs(ctx context.Context, kind string, ids []string, dst interface{}) error {
	keys := BuildStringIDKeys(kind, ids)
	return client.GetByKeys(ctx, keys, dst)
}

// GetByName retrieves model from datastore by name
func (client *Client) GetByName(ctx context.Context, kind string, name string, dst interface{}) error {
	return client.GetByKey(ctx, datastore.NameKey(kind, name, nil), dst)
}

// GetByNames retrieves models from datastore by names
func (client *Client) GetByNames(ctx context.Context, kind string, names []string, dst interface{}) error {
	keys := BuildNameKeys(kind, names)
	return client.GetByKeys(ctx, keys, dst)
}

// GetByQuery retrieves model from datastore by datastore query
func (client *Client) GetByQuery(ctx context.Context, q *datastore.Query, dst interface{}) error {
	keys, err := client.GetAll(ctx, q.KeysOnly(), nil)
	if err != nil {
		return err
	}
	return client.GetByKeys(ctx, keys, dst)
}
