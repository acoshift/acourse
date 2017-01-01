package app

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config type
type Config struct {
	ProjectID string `yaml:"projectId"`
	Email     struct {
		From     string `yaml:"from"`
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"email"`
}

// LoadConfig loads config from file
func LoadConfig(filename string) (*Config, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	err = yaml.Unmarshal(bs, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
