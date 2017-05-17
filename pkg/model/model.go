package model

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func init() {
	gob.Register(&CourseContent{})
}

type scanFunc func(...interface{}) error

// Errors
var (
	ErrNotFound = errors.New("not found")
)

const prefix = "acr:"

// key generate key from given data
func key(xs ...interface{}) string {
	r := prefix
	for _, x := range xs {
		r += ":"
		switch t := x.(type) {
		case string:
			r += t
		case int:
			r += strconv.Itoa(t)
		case int64:
			r += strconv.FormatInt(t, 10)
		case fmt.Stringer:
			r += t.String()
		default:
			panic("key: params not support")
		}
	}
	return r
}

func dec(b []byte, x interface{}) error {
	return gob.NewDecoder(bytes.NewReader(b)).Decode(x)
}

func enc(x interface{}) []byte {
	b := &bytes.Buffer{}
	err := gob.NewEncoder(b).Encode(x)
	if err != nil {
		log.Println("enc:", err)
	}
	return b.Bytes()
}

func mustStmt(stmt *sql.Stmt, err error) *sql.Stmt {
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}
