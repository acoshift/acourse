package webstatic

import (
	"net/http"
	"os"
)

// New creates new webstatic handler
func New(dir string) http.Handler {
	return http.FileServer(&webstaticFS{http.Dir(dir)})
}

type webstaticFS struct {
	http.FileSystem
}

func (fs *webstaticFS) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}
