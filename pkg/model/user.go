package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// User model
type User struct {
	ID        string
	Role      UserRole
	Username  string
	Name      string
	Email     string
	AboutMe   string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRole type
type UserRole struct {
	Admin      bool
	Instructor bool
}

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
	`

	querySaveUser = `
		upsert into users
			(id, name, username, about_me, image, updated_at)
		values
			($1, $2, $3, $4, $5, now())
	`
)

// Save saves user
func (x *User) Save() error {
	if len(x.ID) == 0 {
		return fmt.Errorf("invalid id")
	}
	_, err := db.Exec(querySaveUser, x.ID, x.Name, x.Username, x.AboutMe, x.Image)
	if err != nil {
		return err
	}
	return nil
}

func scanUser(scan scanFunc, x *User) error {
	var admin, instructor sql.NullBool
	var email sql.NullString
	err := scan(&x.ID, &x.Name, &x.Username, &email, &x.AboutMe, &x.Image, &x.CreatedAt, &x.UpdatedAt, &admin, &instructor)
	if err != nil {
		return err
	}
	x.Email = email.String
	x.Role.Admin = admin.Bool
	x.Role.Instructor = instructor.Bool
	return nil
}

// GetUsers gets users
func GetUsers(userIDs []string) ([]*User, error) {
	xs := make([]*User, 0, len(userIDs))
	rows, err := db.Query(queryGetUsers, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x User
		err = scanUser(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}

// GetUser gets user from id
func GetUser(userID string) (*User, error) {
	var x User
	err := scanUser(db.QueryRow(queryGetUser, userID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetUserFromUsername gets user from username
func GetUserFromUsername(username string) (*User, error) {
	var x User
	err := scanUser(db.QueryRow(queryGetUserFromUsername, username).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListUsers lists users
// TODO: pagination
func ListUsers() ([]*User, error) {
	xs := make([]*User, 0)
	rows, err := db.Query(queryListUsers)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var x User
		err = scanUser(rows.Scan, &x)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	return xs, nil
}
