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

// IsUniqueViolation checks is error unique_violation with given constraint,
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
