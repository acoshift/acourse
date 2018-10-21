package course

import (
	"context"
	"database/sql"

	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/course"
)

// Init inits course service
func Init() {
	dispatcher.Register(setOption)
	dispatcher.Register(setImage)
	dispatcher.Register(createContent)
	dispatcher.Register(updateContent)
	dispatcher.Register(getContent)
	dispatcher.Register(deleteContent)
	dispatcher.Register(listContents)
}

func setOption(ctx context.Context, m *course.SetOption) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into course_options
			(course_id, public, enroll, attend, assignment, discount)
		values
			($1, $2, $3, $4, $5, $6)
		on conflict (course_id) do update set
			public = excluded.public,
			enroll = excluded.enroll,
			attend = excluded.attend,
			assignment = excluded.assignment,
			discount = excluded.discount
	`, m.ID, m.Option.Public, m.Option.Enroll, m.Option.Attend, m.Option.Assignment, m.Option.Discount)
	return err
}

func setImage(ctx context.Context, m *course.SetImage) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`update courses set image = $2 where id = $1`, m.ID, m.Image)
	return err
}

func createContent(ctx context.Context, m *course.CreateContent) error {
	// TODO: validate instructor

	q := sqlctx.GetQueryer(ctx)

	return q.QueryRow(`
		insert into course_contents
			(
				course_id,
				i,
				title, long_desc, video_id, video_type
			)
		values
			(
				$1,
				(select coalesce(max(i)+1, 0) from course_contents where course_id = $1),
				$2, $3, $4, $5
			)
		returning id
	`,
		m.ID,
		m.Title, m.LongDesc, m.VideoID, m.VideoType,
	).Scan(&m.Result)
}

func updateContent(ctx context.Context, m *course.UpdateContent) error {
	// TODO: validate ownership

	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update course_contents
		set
			title = $2,
			long_desc = $3,
			video_id = $4,
			updated_at = now()
		where id = $1
	`, m.ContentID, m.Title, m.Desc, m.VideoID)
	return err
}

func getContent(ctx context.Context, m *course.GetContent) error {
	// TODO: validate ownership

	q := sqlctx.GetQueryer(ctx)

	var x course.Content
	err := q.QueryRow(`
		select
			id, course_id, title, long_desc, video_id, video_type, download_url
		from course_contents
		where id = $1
	`, m.ContentID).Scan(
		&x.ID, &x.CourseID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL,
	)
	if err == sql.ErrNoRows {
		return entity.ErrNotFound
	}
	if err != nil {
		return err
	}
	m.Result = &x
	return nil
}

func deleteContent(ctx context.Context, m *course.DeleteContent) error {
	// TODO: validate ownership
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`delete from course_contents where id = $1`, m.ContentID)
	return err
}

func listContents(ctx context.Context, m *course.ListContents) error {
	// TODO: validate ownership

	q := sqlctx.GetQueryer(ctx)

	rows, err := q.Query(`
		select
			id, course_id, title, long_desc, video_id, video_type, download_url
		from course_contents
		where course_id = $1
		order by i
	`, m.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var xs []*course.Content
	for rows.Next() {
		var x course.Content
		err = rows.Scan(
			&x.ID, &x.CourseID, &x.Title, &x.Desc, &x.VideoID, &x.VideoType, &x.DownloadURL,
		)
		if err != nil {
			return err
		}
		xs = append(xs, &x)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	m.Result = xs
	return nil
}
