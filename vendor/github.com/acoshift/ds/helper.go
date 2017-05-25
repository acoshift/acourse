package ds

import (
	"reflect"
)

func valueOf(src interface{}) reflect.Value {
	xs := reflect.ValueOf(src)
	if xs.Kind() == reflect.Ptr {
		xs = xs.Elem()
	}
	return xs
}
