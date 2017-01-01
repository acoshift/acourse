package ctrl

import (
	"github.com/acoshift/acourse/pkg/app"
)

// HealthController implements HealthController interface
type HealthController struct {
}

// NewHealthController creates new controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health runs health action
func (c *HealthController) Health(ctx *app.HealthHealthContext) error {
	return ctx.OK()
}
