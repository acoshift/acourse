package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/acoshift/acourse/pkg/appctx"
	"github.com/acoshift/acourse/pkg/entity"
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

// SaveUser saves user
func SaveUser(ctx context.Context, x *entity.User) error {
	db := appctx.GetDatabase(ctx)

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

func scanUser(scan scanFunc, x *entity.User) error {
	err := scan(&x.ID, &x.Name, &x.Username, &x.Email, &x.AboutMe, &x.Image, &x.CreatedAt, &x.UpdatedAt, &x.Role.Admin, &x.Role.Instructor)
	if err != nil {
		return err
	}
	return nil
}

// GetUsers gets users
func GetUsers(ctx context.Context, userIDs []string) ([]*entity.User, error) {
	db := appctx.GetDatabase(ctx)

	xs := make([]*entity.User, 0, len(userIDs))
	rows, err := db.QueryContext(ctx, queryGetUsers, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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

// GetUser gets user from id
func GetUser(ctx context.Context, userID string) (*entity.User, error) {
	db := appctx.GetDatabase(ctx)

	var x entity.User
	err := scanUser(db.QueryRowContext(ctx, queryGetUser, userID).Scan, &x)
	if err == sql.ErrNoRows {
		return nil, appctx.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetUserFromUsername gets user from username
func GetUserFromUsername(ctx context.Context, username string) (*entity.User, error) {
	db := appctx.GetDatabase(ctx)

	var x entity.User
	err := scanUser(db.QueryRowContext(ctx, queryGetUserFromUsername, username).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := appctx.GetDatabase(ctx)

	var x entity.User
	err := scanUser(db.QueryRowContext(ctx, queryGetUserFromEmail, email).Scan, &x)
	if err == sql.ErrNoRows {
		return nil, appctx.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListUsers lists users
func ListUsers(ctx context.Context, limit, offset int64) ([]*entity.User, error) {
	db := appctx.GetDatabase(ctx)

	xs := make([]*entity.User, 0)
	rows, err := db.QueryContext(ctx, queryListUsers, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
func CountUsers(ctx context.Context) (int64, error) {
	db := appctx.GetDatabase(ctx)

	var cnt int64
	err := db.QueryRowContext(ctx, `select count(*) from users`).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// IsUserExists checks is user exists
func IsUserExists(ctx context.Context, id string) (bool, error) {
	db := appctx.GetDatabase(ctx)

	var cnt int64
	err := db.QueryRowContext(ctx, `select count(*) from users where id = $1`, id).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// CreateUser creates new users
func CreateUser(ctx context.Context, x *entity.User) error {
	db := appctx.GetDatabase(ctx)

	_, err := db.ExecContext(ctx,
		`insert into users (id, username, name, email) values ($1, $2, $3, $4)`,
		x.ID, x.ID, "", x.Email,
	)
	if err != nil {
		return err
	}
	return nil
}
