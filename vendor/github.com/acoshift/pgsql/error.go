package pgsql

import (
	"github.com/lib/pq"
)

func contains(xs []string, x string) bool {
	for _, p := range xs {
		if p == x {
			return true
		}
	}
	return false
}

// IsUniqueViolation checks is error an unique_violation with given constraint,
// constraint can be empty to ignore constraint name checks
func IsUniqueViolation(err error, constraint ...string) bool {
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		if len(constraint) == 0 {
			return true
		}
		return contains(constraint, pqErr.Constraint)
	}
	return false
}

// IsInvalidTextRepresentation checks is error an invalid_text_representation
func IsInvalidTextRepresentation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "22P02" {
		return true
	}
	return false
}
