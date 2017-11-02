package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/acoshift/acourse/pkg/app"
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

	queryListUsers = selectUsers + `
		order by users.created_at desc
		limit $1 offset $2
	`
)

// SaveUser saves user
func SaveUser(ctx context.Context, db DB, x *app.User) error {
	if len(x.ID) == 0 {
		return fmt.Errorf("invalid id")
	}
	_, err := db.ExecContext(ctx, `
		upsert into users
			(id, name, username, about_me, image, updated_at)
		values
			($1, $2, $3, $4, $5, now())
	`, x.ID, x.Name, x.Username, x.AboutMe, x.Image)
	if err != nil {
		return err
	}
	return nil
}

func scanUser(scan scanFunc, x *app.User) error {
	err := scan(&x.ID, &x.Name, &x.Username, &x.Email, &x.AboutMe, &x.Image, &x.CreatedAt, &x.UpdatedAt, &x.Role.Admin, &x.Role.Instructor)
	if err != nil {
		return err
	}
	return nil
}

// GetUsers gets users
func GetUsers(ctx context.Context, db DB, userIDs []string) ([]*app.User, error) {
	xs := make([]*app.User, 0, len(userIDs))
	rows, err := db.QueryContext(ctx, queryGetUsers, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x app.User
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

// GetUser gets user from id
func GetUser(ctx context.Context, db DB, userID string) (*User, error) {
	var x User
	err := scanUser(db.QueryRowContext(ctx, queryGetUser, userID).Scan, &x)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetUserFromUsername gets user from username
func GetUserFromUsername(ctx context.Context, db DB, username string) (*User, error) {
	var x User
	err := scanUser(db.QueryRowContext(ctx, queryGetUserFromUsername, username).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListUsers lists users
func ListUsers(ctx context.Context, db DB, limit, offset int64) ([]*User, error) {
	xs := make([]*User, 0)
	rows, err := db.QueryContext(ctx, queryListUsers, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x User
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
func CountUsers(ctx context.Context, db DB) (int64, error) {
	var cnt int64
	err := db.QueryRowContext(ctx, `select count(*) from users`).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
