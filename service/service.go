package service

import (
	"time"

	"github.com/moonrhythm/dispatcher"
)

// Config is service config
type Config struct {
	Repository     Repository
	BaseURL        string
	Location       *time.Location
	OpenIDCallback string
}

// Init registers service into dispatcher
func Init(cfg Config) {
	s := &svc{cfg}

	dispatcher.Register(s.signUp)
	dispatcher.Register(s.sendPasswordResetEmail)
	dispatcher.Register(s.signInPassword)
	dispatcher.Register(s.generateOpenIDURI)
	dispatcher.Register(s.signInOpenIDCallback)
	dispatcher.Register(s.updateProfile)
	dispatcher.Register(s.createCourse)
	dispatcher.Register(s.updateCourse)
	dispatcher.Register(s.enrollCourse)
	dispatcher.Register(s.createCourseContent)
	dispatcher.Register(s.updateCourseContent)
	dispatcher.Register(s.getCourseContent)
	dispatcher.Register(s.listCourseContents)
	dispatcher.Register(s.deleteCourseContent)
	dispatcher.Register(s.acceptPayment)
	dispatcher.Register(s.rejectPayment)
}

type svc struct {
	Config
}
