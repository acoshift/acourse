package model

import "database/sql"

// model shared vars
var (
	db *sql.DB
)

// Config use to init model package
type Config struct {
	DB *sql.DB
}

// Init inits model package
func Init(config Config) error {
	db = config.DB

	return nil
}
