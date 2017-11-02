package repository

import (
	"github.com/acoshift/acourse/pkg/app"
)

// New creates new repository
func New() app.Repository {
	return &repo{}
}

type repo struct{}

type scanFunc func(...interface{}) error
