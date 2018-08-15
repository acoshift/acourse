package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// NullString scans null into empty string
func NullString(s *string) interface {
	driver.Valuer
	sql.Scanner
} {
	return &nullString{s}
}

type nullString struct {
	value *string
}

func (s *nullString) Scan(src interface{}) error {
	var t sql.NullString
	err := t.Scan(src)
	*s.value = t.String
	return err
}

func (s nullString) Value() (driver.Value, error) {
	if s.value == nil || *s.value == "" {
		return nil, nil
	}
	return *s.value, nil
}
