package session

import (
	"bufio"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeWriter struct{}

func (*fakeWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (*fakeWriter) Header() http.Header {
	return nil
}

func (*fakeWriter) WriteHeader(int) {}

func (*fakeWriter) Push(target string, opts *http.PushOptions) error {
	return nil
}

func (*fakeWriter) Flush() {}

func (*fakeWriter) CloseNotify() <-chan bool {
	return nil
}

func (*fakeWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func TestWriter(t *testing.T) {
	// empty response writer
	w := sessionWriter{}
	w.Push("", nil)
	w.Flush()
	w.CloseNotify()
	w.Hijack()

	w.ResponseWriter = &fakeWriter{}
	w.Push("", nil)
	w.Flush()
	w.CloseNotify()
	w.Hijack()

	called := 0
	w.beforeWriteHeader = func() {
		called++
	}
	w.WriteHeader(200)
	w.WriteHeader(200)
	w.Write([]byte("ok"))

	assert.Equal(t, 1, called)
}
