package service

import (
	"time"

	"github.com/moonrhythm/dispatcher"
)

// Config is service config
type Config struct {
	Repository Repository
	Location   *time.Location
}

// Init registers service into dispatcher
func Init(cfg Config) {
	s := &svc{cfg}

	dispatcher.Register(s.createCourse)
	dispatcher.Register(s.updateCourse)
	dispatcher.Register(s.enrollCourse)
	dispatcher.Register(s.acceptPayment)
	dispatcher.Register(s.rejectPayment)
}

type svc struct {
	Config
}
