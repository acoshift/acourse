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

func (editorRepo) GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.Course
	err := q.QueryRow(`
		select
			id, user_id, title, short_desc, long_desc, image,
			start, url, type, price, courses.discount, enroll_detail,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from courses
			left join course_options as opt on opt.course_id = courses.id
		where id = $1
	`, courseID).Scan(
		&x.ID, &x.UserID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image,
		&x.Start, pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
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
