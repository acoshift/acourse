package ds

import (
	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// NotFound checks is error means not found
func NotFound(err error) bool {
	return err == iterator.Done || err == datastore.ErrNoSuchEntity || err == datastore.ErrInvalidKey
}

// FieldMismatch checks is error field mismatch
func FieldMismatch(err error) bool {
	_, ok := err.(*datastore.ErrFieldMismatch)
	return ok
}

// Ignore removes error(s) from err if f(err)
func Ignore(err error, f func(error) bool) error {
	if f(err) {
		return nil
	}

	if errs, ok := err.(datastore.MultiError); ok {
		if len(errs) == 0 {
			return nil
		}
		es := make(datastore.MultiError, 0)
		for _, err := range errs {
			if !f(err) && err != nil {
				es = append(es, err)
			}
		}
		if len(es) > 0 {
			return es
		}
		return nil
	}
	return err
}

// IgnoreFieldMismatch returns nil if err is field mismatch error(s)
func IgnoreFieldMismatch(err error) error {
	return Ignore(err, FieldMismatch)
}

// IgnoreNotFound returns nil if err is not found error(s)
func IgnoreNotFound(err error) error {
	return Ignore(err, NotFound)
}
