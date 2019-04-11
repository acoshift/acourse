package config

import (
	"github.com/acoshift/configfile"
)

var cfg = configfile.NewReader("config")

var (
	StringDefault   = cfg.StringDefault
	String          = cfg.String
	IntDefault      = cfg.IntDefault
	Int             = cfg.Int
	DurationDefault = cfg.DurationDefault
	Bytes           = cfg.Bytes
)
