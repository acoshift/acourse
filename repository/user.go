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
