package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
)

const (
	selectUsers = `
		select
			users.id,
			users.name,
			users.username,
			users.email,
			users.about_me,
			users.image,
			users.created_at,
			users.updated_at,
			roles.admin,
			roles.instructor
		from users
			left join roles on users.id = roles.user_id
	`

	queryGetUsers = selectUsers + `
		where users.id = any($1)
	`

	queryGetUser = selectUsers + `
		where users.id = $1
	`

	queryGetUserFromUsername = selectUsers + `
		where users.username = $1
	`

	queryGetUserFromEmail = selectUsers + `
		where users.email = $1
	`

	queryListUsers = selectUsers + `
		order by users.created_at desc
		limit $1 offset $2
	`
)

func scanUser(scan scanFunc, x *entity.User) error {
	err := scan(&x.ID, &x.Name, &x.Username, &x.Email, &x.AboutMe, &x.Image, &x.CreatedAt, &x.UpdatedAt, &x.Role.Admin, &x.Role.Instructor)
	if err != nil {
		return err
	}
	return nil
}

// GetUser gets user from id
func GetUser(ctx context.Context, userID string) (*entity.User, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.User
	err := scanUser(q.QueryRow(queryGetUser, userID).Scan, &x)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetUserFromUsername gets user from username
func GetUserFromUsername(ctx context.Context, username string) (*entity.User, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.User
	err := scanUser(q.QueryRow(queryGetUserFromUsername, username).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// FindUserByEmail finds user by email
func FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.User
	err := scanUser(q.QueryRow(queryGetUserFromEmail, email).Scan, &x)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListUsers lists users
func ListUsers(ctx context.Context, limit, offset int64) ([]*entity.User, error) {
	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(queryListUsers, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*entity.User
	for rows.Next() {
		var x entity.User
		err = scanUser(rows.Scan, &x)
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
func RegisterUser(ctx context.Context, x *entity.User) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into users
			(id, username, name, email, image)
		values
			($1, $2, $3, $4, $5)
	`, x.ID, x.Username, x.Name, x.Email, x.Image)
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
func UpdateUser(ctx context.Context, userID string, username, name, aboutMe string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update users
		set
			username = $2,
			name = $3,
			about_me = $4,
			updated_at = now()
		where id = $1
	`, userID, username, name, aboutMe)
	return err
}
