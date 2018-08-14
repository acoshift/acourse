// +build cgo,cbrotli

package middleware

import (
	"io"

	"github.com/google/brotli/go/cbrotli"
)

// BrCompressor is the brotli compressor for compress middleware
var BrCompressor = CompressConfig{
	Skipper: DefaultSkipper,
	New: func() Compressor {
		return &brWriter{quality: 4}
	},
	Encoding:  "br",
	Vary:      defaultCompressVary,
	Types:     defaultCompressTypes,
	MinLength: defaultCompressMinLength,
}

type brWriter struct {
	quality int
	*cbrotli.Writer
}

func (w *brWriter) Reset(p io.Writer) {
	w.Writer = cbrotli.NewWriter(p, cbrotli.WriterOptions{Quality: w.quality})
}
