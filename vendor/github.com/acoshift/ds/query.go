package ds

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

// Query is the function for set datastore query
type Query func(q *datastore.Query) *datastore.Query

// Query run Get All
// dst is *[]*Model
func (client *Client) Query(ctx context.Context, kind string, dst interface{}, qs ...Query) error {
	q := datastore.NewQuery(kind)
	for _, setter := range qs {
		q = setter(q)
	}
	q = q.KeysOnly()

	keys, err := client.GetAll(ctx, q, nil)
	if err != nil {
		return err
	}
	return client.GetByKeys(ctx, keys, dst)
}

// QueryFirst run Get to get the first result
func (client *Client) QueryFirst(ctx context.Context, kind string, dst interface{}, qs ...Query) error {
	q := datastore.NewQuery(kind)
	for _, setter := range qs {
		q = setter(q)
	}
	q = q.Limit(1).KeysOnly()

	key, err := client.Run(ctx, q).Next(nil)
	if err != nil {
		return err
	}
	return client.GetByKey(ctx, key, dst)
}

// QueryKeys queries only key
func (client *Client) QueryKeys(ctx context.Context, kind string, qs ...Query) ([]*datastore.Key, error) {
	q := datastore.NewQuery(kind)
	for _, setter := range qs {
		q = setter(q)
	}
	q = q.KeysOnly()

	keys, err := client.GetAll(ctx, q, nil)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// QueryCount counts entity
func (client *Client) QueryCount(ctx context.Context, kind string, qs ...Query) (int, error) {
	q := datastore.NewQuery(kind)
	for _, setter := range qs {
		q = setter(q)
	}

	return client.Count(ctx, q)
}

// Query Helper functions

// Filter func
func Filter(filterStr string, value interface{}) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Filter(filterStr, value)
	}
}

// CreateBefore quries is model created after (or equals) given time
func CreateBefore(t time.Time, equals bool) Query {
	p := "CreatedAt <"
	if equals {
		p += "="
	}
	return Filter(p, t)
}

// CreateAfter quries is model created after (or equals) given time
func CreateAfter(t time.Time, equals bool) Query {
	p := "CreatedAt >"
	if equals {
		p += "="
	}
	return Filter(p, t)
}

// UpdateBefore queries is model updated before (or equals) given time
func UpdateBefore(t time.Time, equals bool) Query {
	p := "UpdatedAt <"
	if equals {
		p += "="
	}
	return Filter(p, t)
}

// UpdateAfter queries is model updated after (or equals) given time
func UpdateAfter(t time.Time, equals bool) Query {
	p := "UpdatedAt >"
	if equals {
		p += "="
	}
	return Filter(p, t)
}

// Offset adds offset to query
func Offset(offset int) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Offset(offset)
	}
}

// Limit adds limit to query
func Limit(limit int) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Limit(limit)
	}
}

// Namespace adds namespace to query
func Namespace(ns string) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Namespace(ns)
	}
}

// Order adds order to query
func Order(fieldName string) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Order(fieldName)
	}
}

// Project adds order to query
func Project(fieldNames ...string) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Project(fieldNames...)
	}
}

// Transaction adds transaction to query
func Transaction(t *Tx) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Transaction(t.Transaction)
	}
}

// Ancestor adds ancestor to query
func Ancestor(ancestor *datastore.Key) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Ancestor(ancestor)
	}
}

// EventualConsistency adds eventual consistency to query
func EventualConsistency() Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.EventualConsistency()
	}
}

// Distinct adds distinct to query
func Distinct() Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Distinct()
	}
}

// DistinctOn adds distinct on to query
func DistinctOn(fieldNames ...string) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.DistinctOn(fieldNames...)
	}
}
