package configfile

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/configfile/internal/reader"
)

// NewReader creates new config reader
func NewReader(base string) *Reader {
	stats, _ := os.Stat(base)
	if stats != nil && !stats.IsDir() {
		return &Reader{reader.NewYAML(base)}
	}
	return &Reader{reader.NewDir(base)}
}

// NewDirReader creates new config dir reader
func NewDirReader(base string) *Reader {
	return &Reader{reader.NewDir(base)}
}

// NewYAMLReader creates new yaml reader
func NewYAMLReader(filename string) *Reader {
	return &Reader{reader.NewYAML(filename)}
}

type intlReader interface {
	Read(name string) ([]byte, error)
}

// Reader is the config reader
type Reader struct {
	r intlReader
}

func (r *Reader) readString(name string) (string, error) {
	b, err := r.r.Read(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *Reader) readInt(name string) (int, error) {
	s, err := r.readString(name)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func (r *Reader) readInt64(name string) (int64, error) {
	s, err := r.readString(name)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (r *Reader) readBool(name string) (bool, error) {
	s, err := r.readString(name)
	if err != nil {
		return false, err
	}
	if s == "" {
		return false, io.EOF
	}
	if s == "0" {
		return false, nil
	}
	if strings.ToLower(s) == "false" {
		return false, nil
	}
	return true, nil
}

func (r *Reader) readDuration(name string) (time.Duration, error) {
	s, err := r.readString(name)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(s)
}

// BytesDefault reads bytes from config file with default value
func (r *Reader) BytesDefault(name string, def []byte) []byte {
	b, err := r.r.Read(name)
	if err != nil {
		return def
	}
	return b
}

// Bytes reads bytes from config file
func (r *Reader) Bytes(name string) []byte {
	return r.BytesDefault(name, []byte{})
}

// MustBytes reads bytes from config file, panic if file not exists
func (r *Reader) MustBytes(name string) []byte {
	s, err := r.r.Read(name)
	if err != nil {
		panic(err)
	}
	return s
}

// StringDefault reads string from config file with default value
func (r *Reader) StringDefault(name string, def string) string {
	s, err := r.readString(name)
	if err != nil {
		return def
	}
	return s
}

// String reads string from config file
func (r *Reader) String(name string) string {
	return r.StringDefault(name, "")
}

// MustString reads string from config file, panic if file not exists
func (r *Reader) MustString(name string) string {
	s, err := r.readString(name)
	if err != nil {
		panic(err)
	}
	return s
}

// IntDefault reads int from config file with default value
func (r *Reader) IntDefault(name string, def int) int {
	i, err := r.readInt(name)
	if err != nil {
		return def
	}
	return i
}

// Int reads int from config file
func (r *Reader) Int(name string) int {
	return r.IntDefault(name, 0)
}

// MustInt reads int from config file, panic if file not exists or data can not parse to int
func (r *Reader) MustInt(name string) int {
	i, err := r.readInt(name)
	if err != nil {
		panic(err)
	}
	return i
}

// Int64Default reads int64 from config file with default value
func (r *Reader) Int64Default(name string, def int64) int64 {
	i, err := r.readInt64(name)
	if err != nil {
		return def
	}
	return i
}

// Int64 reads int from config file
func (r *Reader) Int64(name string) int64 {
	return r.Int64Default(name, 0)
}

// MustInt64 reads int64 from config file, panic if file not exists or data can not parse to int64
func (r *Reader) MustInt64(name string) int64 {
	i, err := r.readInt64(name)
	if err != nil {
		panic(err)
	}
	return i
}

// BoolDefault reads bool from config file with default value,
// result is false if lower case data is "", "0", or "false", otherwise true
func (r *Reader) BoolDefault(name string, def bool) bool {
	b, err := r.readBool(name)
	if err != nil {
		return def
	}
	return b
}

// Bool reads bool from config file, see BoolDefault
func (r *Reader) Bool(name string) bool {
	return r.BoolDefault(name, false)
}

// MustBool reads bool from config file, see BoolDefault,
// panic if file not exists
func (r *Reader) MustBool(name string) bool {
	b, err := r.readBool(name)
	if err != nil {
		panic(err)
	}
	return b
}

// DurationDefault reads string then parse as duration from config file with default value
func (r *Reader) DurationDefault(name string, def time.Duration) time.Duration {
	d, err := r.readDuration(name)
	if err != nil {
		return def
	}
	return d
}

// Duration reads string then parse as duration from config file
func (r *Reader) Duration(name string) time.Duration {
	return r.DurationDefault(name, 0)
}

// MustDuration reads string then parse as duration from config file,
// panic if file not exists
func (r *Reader) MustDuration(name string) time.Duration {
	b, err := r.readDuration(name)
	if err != nil {
		panic(err)
	}
	return b
}
