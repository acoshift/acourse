package repository

import (
	"context"

	"github.com/lib/pq"

	"github.com/acoshift/acourse/pkg/app"
)

func (repo) FindUsers(ctx context.Context, userIDs []string) ([]*app.User, error) {
	db := app.GetDatabase(ctx)

	xs := make([]*app.User, 0, len(userIDs))
	rows, err := db.QueryContext(ctx, `
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
		where users.id = $1
	`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x app.User
		err = rows.Scan(&x.ID, &x.Name, &x.Username, &x.Email, &x.AboutMe, &x.Image, &x.CreatedAt, &x.UpdatedAt, &x.Role.Admin, &x.Role.Instructor)
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
