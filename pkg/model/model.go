package model

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
)

type scanFunc func(...interface{}) error

// Errors
var (
	ErrNotFound = errors.New("not found")
)

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
