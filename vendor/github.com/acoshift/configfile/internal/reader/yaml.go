package reader

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

// NewYAML creates new yaml reader
func NewYAML(filename string) *YAML {
	var r YAML
	fs, _ := os.Open(filename)
	if fs != nil {
		yaml.NewDecoder(fs).Decode(&r.d)
	}
	return &r
}

// YAML reads config from yaml file
type YAML struct {
	d map[string]string
}

// Read reads a config
func (r *YAML) Read(name string) ([]byte, error) {
	p, ok := r.d[name]
	if !ok {
		return nil, errNotFound
	}
	return []byte(p), nil
}
