package health

import (
	"context"

	"github.com/acoshift/acourse/pkg/app"
)

// New creates new health service
func New() app.HealthService {
	return &service{}
}

type service struct{}

func (*service) Check(ctx context.Context) error {
	return nil
}
