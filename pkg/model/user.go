package model

import (
	"fmt"
	"time"

	"github.com/acoshift/acourse/pkg/internal"
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

const selectUsers = `
	SELECT
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
	FROM users
		LEFT JOIN roles ON users.id = roles.id
`

var (
	getUsersStmt, _ = internal.GetDB().Prepare(selectUsers + `
		WHERE users.id = ANY($1);
	`)

	getUserStmt, _ = internal.GetDB().Prepare(selectUsers + `
		WHERE users.id = $1;
	`)

	getUserFromUsernameStmt, _ = internal.GetDB().Prepare(selectUsers + `
		WHERE users.username = $1;
	`)

	listUsersStmt, _ = internal.GetDB().Prepare(selectUsers + `
		ORDER BY users.created_at DESC;
	`)

	saveUserStmt, _ = internal.GetDB().Prepare(`
		UPSERT INTO users
			(id, name, username, about_me, image, updated_at)
		VALUES
			($1, $2, $3, $4, $5, now());
	`)
)

// Save saves user
func (x *User) Save() error {
	if len(x.ID) == 0 {
		return fmt.Errorf("invalid id")
	}
	_, err := saveUserStmt.Exec(x.ID, x.Name, x.Username, x.AboutMe, x.Image)
	if err != nil {
		return err
	}
	return nil
}

func scanUser(scan scanFunc, x *User) error {
	var admin, instructor *bool
	var email *string
	err := scan(&x.ID, &x.Name, &x.Username, &email, &x.AboutMe, &x.Image, &x.CreatedAt, &x.UpdatedAt, &admin, &instructor)
	if err != nil {
		return err
	}
	if email != nil {
		x.Email = *email
	}
	if admin != nil {
		x.Role.Admin = *admin
	}
	if instructor != nil {
		x.Role.Instructor = *instructor
	}
	return nil
}

// GetUsers gets users
func GetUsers(userIDs []string) ([]*User, error) {
	xs := make([]*User, 0, len(userIDs))
	rows, err := getUsersStmt.Query(pq.Array(userIDs))
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
	err := scanUser(getUserStmt.QueryRow(userID).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// GetUserFromUsername gets user from username
func GetUserFromUsername(username string) (*User, error) {
	var x User
	err := scanUser(getUserFromUsernameStmt.QueryRow(username).Scan, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

// ListUsers lists users
// TODO: pagination
func ListUsers() ([]*User, error) {
	xs := make([]*User, 0)
	rows, err := listUsersStmt.Query()
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
