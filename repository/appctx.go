package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/acourse/context/appctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/user"
)

// NewAppCtx creates new appctx repository
func NewAppCtx() appctx.Repository {
	return &appctxRepo{}
}

type appctxRepo struct {
}

func (appctxRepo) GetUser(ctx context.Context, userID string) (*user.User, error) {
	q := sqlctx.GetQueryer(ctx)

	var x user.User
	err := q.QueryRow(`
		select
			u.id, u.name, u.username, coalesce(u.email, ''), u.about_me, u.image,
			coalesce(r.admin, false), coalesce(r.instructor, false)
		from users as u
			left join roles as r on u.id = r.user_id
		where u.id = $1
	`, userID).Scan(
		&x.ID, &x.Name, &x.Username, &x.Email, &x.AboutMe, &x.Image,
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
