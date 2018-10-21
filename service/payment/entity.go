package payment

import (
	"time"

	"github.com/lib/pq"
)

// Payment type
type Payment struct {
	ID            string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
	CreatedAt     time.Time
	At            pq.NullTime

	User struct {
		ID       string
		Username string
		Name     string
		Email    string
		Image    string
	}
	Course struct {
		ID    string
		Title string
		Image string
		URL   string
	}
}
