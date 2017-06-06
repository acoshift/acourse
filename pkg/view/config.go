package view

// view shared vars
var (
	baseURL string
)

// Config use to init view package
type Config struct {
	BaseURL string
}

// Init inits view package
func Init(config Config) error {
	baseURL = config.BaseURL

	return nil
}
