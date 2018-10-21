package file

import (
	"io"

	"github.com/satori/go.uuid"
)

// GenerateFilename generates new filename
func GenerateFilename() string {
	return "upload/" + uuid.NewV4().String()
}

// Store stores file
type Store struct {
	Reader   io.Reader
	Filename string
	Async    bool

	Result string // download url
}
