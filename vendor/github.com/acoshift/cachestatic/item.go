package cachestatic

import (
	"net/http"
	"time"

	"github.com/acoshift/header"
)

type item struct {
	data    []byte
	header  http.Header
	modTime time.Time
}

func createItem(w *responseWriter) *item {
	it := item{
		data:   w.cache.Bytes(),
		header: w.h,
	}
	if w.h != nil {
		if v := w.h.Get(header.LastModified); len(v) > 0 {
			it.modTime, _ = time.Parse(http.TimeFormat, v)
		} else {
			it.modTime = time.Now().UTC().Round(time.Second)
			w.h.Set(header.LastModified, it.modTime.Format(http.TimeFormat))
		}
	}
	return &it
}
