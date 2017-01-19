package store

import (
	"time"

	"cloud.google.com/go/datastore"
)

// Query is the query function
type Query func(*datastore.Query) *datastore.Query

func queryFilter(filterStr string, value interface{}) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Filter(filterStr, value)
	}
}

// QueryCreateBefore queries is model created before (or equals) given time
func QueryCreateBefore(t time.Time, equals bool) Query {
	return func(q *datastore.Query) *datastore.Query {
		p := "CreatedAt <"
		if equals {
			p += "="
		}
		return q.Filter(p, t)
	}
}

// QueryCreateAfter queries is model created after (or equals) given time
func QueryCreateAfter(t time.Time, equals bool) Query {
	return func(q *datastore.Query) *datastore.Query {
		p := "CreatedAt >"
		if equals {
			p += "="
		}
		return q.Filter(p, t)
	}
}

// QueryUpdateBefore queries is model updated before (or equals) given time
func QueryUpdateBefore(t time.Time, equals bool) Query {
	return func(q *datastore.Query) *datastore.Query {
		p := "UpdatedAt <"
		if equals {
			p += "="
		}
		return q.Filter(p, t)
	}
}

// QueryUpdateAfter queries is model updated after (or equals) given time
func QueryUpdateAfter(t time.Time, equals bool) Query {
	return func(q *datastore.Query) *datastore.Query {
		p := "UpdatedAt >"
		if equals {
			p += "="
		}
		return q.Filter(p, t)
	}
}

// QueryArchive queries is model archived
func QueryArchive(archive bool) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Filter("Archive =", archive)
	}
}

// QueryOffset adds offset to query
func QueryOffset(offset int) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Offset(offset)
	}
}

// QueryLimit adds limit to query
func QueryLimit(limit int) Query {
	return func(q *datastore.Query) *datastore.Query {
		return q.Limit(limit)
	}
}

// QueryUserID adds user id to query
func QueryUserID(userID string) Query {
	return queryFilter("UserID =", userID)
}

// QueryAssignmentID adds assignment id to query
func QueryAssignmentID(assignmentID string) Query {
	return queryFilter("AssignmentID =", assignmentID)
}

// QueryCourseID adds course id to query
func QueryCourseID(courseID string) Query {
	return queryFilter("CourseID =", courseID)
}
