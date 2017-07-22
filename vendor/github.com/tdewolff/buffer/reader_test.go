package buffer // import "github.com/tdewolff/buffer"

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	s := []byte("abcde")
	r := NewReader(s)
	assert.Equal(t, s, r.Bytes(), "reader must return bytes stored")

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, 3, n, "first read must read 3 characters")
	assert.Equal(t, []byte("abc"), buf, "first read must match 'abc'")

	n, err = r.Read(buf)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, 2, n, "second read must read 2 characters")
	assert.Equal(t, []byte("de"), buf[:n], "second read must match 'de'")

	n, err = r.Read(buf)
	assert.Equal(t, io.EOF, err, "error must be io.EOF")
	assert.Equal(t, 0, n, "third read must read 0 characters")

	n, err = r.Read(nil)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, 0, n, "read to nil buffer must return 0 characters read")
}

func ExampleNewReader() {
	r := NewReader([]byte("Lorem ipsum"))
	w := &bytes.Buffer{}
	io.Copy(w, r)
	fmt.Println(w.String())
	// Output: Lorem ipsum
}
