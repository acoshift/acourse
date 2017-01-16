package health

import (
	"github.com/acoshift/acourse/pkg/app"
)

// New creates new health controller
func New() app.HealthController {
	return &healthController{}
}

type healthController struct{}

func (*healthController) Check() error {
	return nil
}
