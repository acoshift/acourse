package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/controller/editor"
	"github.com/acoshift/acourse/entity"
)

// NewEditor creates new editor repository
func NewEditor() editor.Repository {
	return &editorRepo{}
}

type editorRepo struct {
}

func (editorRepo) GetCourseUserID(ctx context.Context, courseID string) (userID string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select user_id from courses where id = $1`, courseID).Scan(&userID)
	if err == sql.ErrNoRows {
		err = entity.ErrNotFound
	}
	return
}

func (editorRepo) GetCourseURL(ctx context.Context, courseID string) (url string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`select url from courses where id = $1`, courseID).Scan(pgsql.NullString(&url))
	return
}
