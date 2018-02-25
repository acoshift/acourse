package repository

import (
	"database/sql"

	"github.com/lib/pq"

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

// GetUsers gets users
func GetUsers(q Queryer, userIDs []string) ([]*entity.User, error) {
	xs := make([]*entity.User, 0, len(userIDs))
	rows, err := q.Query(queryGetUsers, pq.Array(userIDs))
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
func GetUser(q Queryer, userID string) (*entity.User, error) {
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
func GetUserFromUsername(q Queryer, username string) (*entity.User, error) {
	var x entity.User
	err := scanUser(q.QueryRow(queryGetUserFromUsername, username).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// FindUserByEmail finds user by email
func FindUserByEmail(q Queryer, email string) (*entity.User, error) {
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
func ListUsers(q Queryer, limit, offset int64) ([]*entity.User, error) {
	xs := make([]*entity.User, 0)
	rows, err := q.Query(queryListUsers, limit, offset)
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
func CountUsers(q Queryer) (int64, error) {
	var cnt int64
	err := q.QueryRow(`
		SELECT count(*)
		  FROM users;
	`).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// IsUserExists checks is user exists
func IsUserExists(q Queryer, id string) (bool, error) {
	var cnt int64
	err := q.QueryRow(`
		SELECT count(*)
		  FROM users
		 WHERE id = $1;
	`, id).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// CreateUser creates new users
func CreateUser(q Queryer, x *entity.User) error {
	_, err := q.Exec(`
		INSERT INTO users (id, username, name, email)
		VALUES ($1, $2, $3, $4);
	`, x.ID, x.ID, "", x.Email)
	if err != nil {
		return err
	}
	return nil
}
