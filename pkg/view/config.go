package view

// view shared vars
var (
	xsrfSecret string
)

// Config use to init view package
type Config struct {
	XSRFSecret string
}

// Init inits view package
func Init(config Config) error {
	xsrfSecret = config.XSRFSecret

	return nil
}
