package ctrl

import (
	"github.com/acoshift/acourse/pkg/app"
)

// NewHealth creates new health controller
func NewHealth() app.HealthController {
	return &healthController{}
}

type healthController struct{}

func (*healthController) Check() error {
	return nil
}
