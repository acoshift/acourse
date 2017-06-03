package view

// view shared vars
var (
	xsrfSecret string
	baseURL    string
)

// Config use to init view package
type Config struct {
	XSRFSecret string
	BaseURL    string
}

// Init inits view package
func Init(config Config) error {
	xsrfSecret = config.XSRFSecret
	baseURL = config.BaseURL

	return nil
}
