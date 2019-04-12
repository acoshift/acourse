package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"

	"github.com/acoshift/acourse/internal/pkg/app"
)

// Errors
var (
	ErrUsernameNotAvailable = errors.New("user: username not available")
	ErrEmailNotAvailable    = errors.New("user: email not available")
)

// User type
type User struct {
	ID       string
	Role     Role
	Username string
	Name     string
	Email    string
	AboutMe  string
	Image    string
}

// Role type
type Role struct {
	Admin      bool
	Instructor bool
}

type CreateArgs struct {
	ID       string
	Username string
	Name     string
	Email    string
	Image    string
}

// Create creates new user
func Create(ctx context.Context, m *CreateArgs) error {
	_, err := pgctx.Exec(ctx, `
		insert into users
			(id, username, name, email, image)
		values
			($1, $2, $3, $4, $5)
	`, m.ID, m.Username, m.Name, pgsql.NullString(&m.Email), m.Image)
	if pgsql.IsUniqueViolation(err, "users_email_key") {
		return ErrEmailNotAvailable
	}
	if pgsql.IsUniqueViolation(err, "users_username_key") {
		return ErrUsernameNotAvailable
	}
	return err
}

type UpdateArgs struct {
	ID       string
	Username string
	Name     string
	AboutMe  string
}

// Update updates user
func Update(ctx context.Context, m *UpdateArgs) error {
	_, err := pgctx.Exec(ctx, `
		update users
		set
			username = $2,
			name = $3,
			about_me = $4,
			updated_at = now()
		where id = $1
	`, m.ID, m.Username, m.Name, m.AboutMe)
	return err
}

// IsExists checks is user exists
func IsExists(ctx context.Context, id string) (bool, error) {
	var b bool
	err := pgctx.QueryRow(ctx, `
		select exists (
			select 1
			from users
			where id = $1
		)
	`, id).Scan(&b)
	return b, err
}

// SetImage sets user image
func SetImage(ctx context.Context, id string, image string) error {
	_, err := pgctx.Exec(ctx, `
		update users
		set image = $2
		where id = $1
	`, id, image)
	return err
}

// Get gets user from id
func Get(ctx context.Context, id string) (*User, error) {
	var x User
	err := pgctx.QueryRow(ctx, `
		select
			u.id, u.name, u.username, coalesce(u.email, ''), u.about_me, u.image,
			coalesce(r.admin, false), coalesce(r.instructor, false)
		from users as u
			left join roles as r on u.id = r.user_id
		where u.id = $1
	`, id).Scan(
		&x.ID, &x.Name, &x.Username, &x.Email, &x.AboutMe, &x.Image,
		&x.Role.Admin, &x.Role.Instructor,
	)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
