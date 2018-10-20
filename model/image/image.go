package image

import (
	"io"
)

// JPEG resize and encode image to JPEG
type JPEG struct {
	Writer  io.Writer
	Reader  io.Reader
	Width   int
	Height  int
	Quality int
	Crop    bool
}
