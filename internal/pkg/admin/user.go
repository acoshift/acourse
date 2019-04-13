package admin

import (
	"context"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"
)

type User struct {
	ID        string
	Username  string
	Name      string
	Email     string
	Image     string
	CreatedAt time.Time
}

func GetUsers(ctx context.Context, limit, offset int64) ([]*User, error) {
	// language=SQL
	rows, err := pgctx.Query(ctx, `
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

	var xs []*User
	for rows.Next() {
		var x User
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

func CountUsers(ctx context.Context) (cnt int64, err error) {
	// language=SQL
	err = pgctx.QueryRow(ctx, `select count(*) from users`).Scan(&cnt)
	return
}
