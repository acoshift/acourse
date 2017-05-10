package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
)

func init() {
	gob.Register(CourseType(0))
	gob.Register(&CourseContent{})
}

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
