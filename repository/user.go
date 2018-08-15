package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
)

// GetUser gets user by id
func GetUser(ctx context.Context, userID string) (*entity.User, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.User
	err := q.QueryRow(`
		select
			users.id,
			users.name,
			users.username,
			users.email,
			users.about_me,
			users.image,
			coalesce(roles.admin, false),
			coalesce(roles.instructor, false)
		from users
			left join roles on users.id = roles.user_id
		where users.id = $1
	`, userID).Scan(
		&x.ID, &x.Name, &x.Username, pgsql.NullString(&x.Email), &x.AboutMe, &x.Image,
		&x.Role.Admin, &x.Role.Instructor,
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetEmailSignInUserByEmail gets email sign in user by email
func GetEmailSignInUserByEmail(ctx context.Context, email string) (*entity.EmailSignInUser, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.EmailSignInUser
	err := q.QueryRow(`
		select
			id, name, email
		from users
		where email = $1
	`, email).Scan(
		&x.ID, &x.Name, pgsql.NullString(&x.Email),
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListUsers lists users
func ListUsers(ctx context.Context, limit, offset int64) ([]*entity.UserItem, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		select
			id, name, username, email,
			image, created_at
		from users
		order by created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*entity.UserItem
	for rows.Next() {
		var x entity.UserItem
		err = rows.Scan(
			&x.ID, &x.Name, &x.Username, pgsql.NullString(&x.Email),
			&x.Image, &x.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return xs, nil
}

// CountUsers counts users
func CountUsers(ctx context.Context) (cnt int64, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select count(*) from users`).Scan(&cnt)
	return
}

// IsUserExists checks is user exists
func IsUserExists(ctx context.Context, id string) (exists bool, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select exists (
			select 1
			from users
			where id = $1
		)
	`, id).Scan(&exists)
	return
}

// RegisterUser registers new users
func RegisterUser(ctx context.Context, x *entity.RegisterUser) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into users
			(id, username, name, email, image)
		values
			($1, $2, $3, $4, $5)
	`, x.ID, x.Username, x.Name, pgsql.NullString(&x.Email), x.Image)
	return err
}

// SetUserImage sets user's image
func SetUserImage(ctx context.Context, userID string, image string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update users
		set image = $2
		where id = $1
	`, userID, image)
	return err
}

// UpdateUser updates user
func UpdateUser(ctx context.Context, x *entity.UpdateUser) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update users
		set
			username = $2,
			name = $3,
			about_me = $4,
			updated_at = now()
		where id = $1
	`, x.ID, x.Username, x.Name, x.AboutMe)
	return err
}
