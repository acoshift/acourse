// +build !cbrotli

package middleware

// BrCompressor is a noop compressor fallback for br
var BrCompressor = CompressConfig{
	Skipper: AlwaysSkip,
}
