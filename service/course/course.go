package course

import (
	"bytes"
	"context"
	"database/sql"
	"io"

	"github.com/acoshift/pgsql"
	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/app"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/model/file"
	"github.com/acoshift/acourse/model/image"
)

// Init inits course service
func Init() {
	dispatcher.Register(setOption)
	dispatcher.Register(setImage)
	dispatcher.Register(getURL)
	dispatcher.Register(getUserID)
	dispatcher.Register(get)
	dispatcher.Register(createContent)
	dispatcher.Register(updateContent)
	dispatcher.Register(getContent)
	dispatcher.Register(deleteContent)
	dispatcher.Register(listContents)
	dispatcher.Register(insertEnroll)
	dispatcher.Register(create)
	dispatcher.Register(update)
}

func setOption(ctx context.Context, m *course.SetOption) error {
	_, err := sqlctx.Exec(ctx, `
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
	_, err := sqlctx.Exec(ctx, `update courses set image = $2 where id = $1`, m.ID, m.Image)
	return err
}

func getURL(ctx context.Context, m *course.GetURL) error {
	return sqlctx.QueryRow(ctx,
		`select url from courses where id = $1`,
		m.ID,
	).Scan(pgsql.NullString(&m.Result))
}

func getUserID(ctx context.Context, m *course.GetUserID) error {
	err := sqlctx.QueryRow(ctx, `select user_id from courses where id = $1`, m.ID).Scan(&m.Result)
	if err == sql.ErrNoRows {
		return entity.ErrNotFound
	}
	return err
}

func get(ctx context.Context, m *course.Get) error {
	var x course.Course
	err := sqlctx.QueryRow(ctx, `
		select
			id, user_id, title, short_desc, long_desc, image,
			start, url, type, price, courses.discount, enroll_detail,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from courses
			left join course_options as opt on opt.course_id = courses.id
		where id = $1
	`, m.ID).Scan(
		&x.ID, &x.UserID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image,
		&x.Start, pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
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

func createContent(ctx context.Context, m *course.CreateContent) error {
	// TODO: validate instructor

	return sqlctx.QueryRow(ctx, `
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

	_, err := sqlctx.Exec(ctx, `
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

	var x course.Content
	err := sqlctx.QueryRow(ctx, `
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

	_, err := sqlctx.Exec(ctx, `delete from course_contents where id = $1`, m.ContentID)
	return err
}

func listContents(ctx context.Context, m *course.ListContents) error {
	// TODO: validate ownership

	rows, err := sqlctx.Query(ctx, `
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

func insertEnroll(ctx context.Context, m *course.InsertEnroll) error {
	_, err := sqlctx.Exec(ctx, `
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, m.UserID, m.ID)
	return err
}

func create(ctx context.Context, m *course.Create) error {
	// TODO: validate user role

	if m.Title == "" {
		return app.NewUIError("title required")
	}

	var imageURL string
	if m.Image != nil {
		err := image.Validate(m.Image)
		if err != nil {
			return err
		}

		image, err := m.Image.Open()
		if err != nil {
			return err
		}
		defer image.Close()

		imageURL, err = uploadCourseCoverImage(ctx, image)
		image.Close()
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	return sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		err := sqlctx.QueryRow(ctx, `
			insert into courses
				(user_id, title, short_desc, long_desc, image, start)
			values
				($1, $2, $3, $4, $5, $6)
			returning id
		`, m.UserID, m.Title, m.ShortDesc, m.LongDesc, imageURL, pgsql.NullTime(&m.Start)).Scan(&m.Result)
		if err != nil {
			return err
		}

		return dispatcher.Dispatch(ctx, &course.SetOption{ID: m.Result, Option: course.Option{}})
	})
}

func uploadCourseCoverImage(ctx context.Context, r io.Reader) (string, error) {
	buf := &bytes.Buffer{}

	if err := dispatcher.Dispatch(ctx, &image.JPEG{
		Writer:  buf,
		Reader:  r,
		Width:   1200,
		Quality: 90,
	}); err != nil {
		return "", err
	}

	filename := file.GenerateFilename() + ".jpg"

	store := file.Store{Reader: buf, Filename: filename}
	if err := dispatcher.Dispatch(ctx, &store); err != nil {
		return "", err
	}
	return store.Result, nil
}

func update(ctx context.Context, m *course.Update) error {
	// TODO: validate user role
	// user := appctx.GetUser(ctx)

	if m.ID == "" {
		return app.NewUIError("course id required")
	}

	if m.Title == "" {
		return app.NewUIError("title required")
	}

	var imageURL string
	if m.Image != nil {
		err := image.Validate(m.Image)
		if err != nil {
			return err
		}

		image, err := m.Image.Open()
		if err != nil {
			return err
		}
		defer image.Close()

		imageURL, err = uploadCourseCoverImage(ctx, image)
		image.Close()
		if err != nil {
			return app.NewUIError(err.Error())
		}
	}

	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		_, err := sqlctx.Exec(ctx, `
			update courses
			set
				title = $2,
				short_desc = $3,
				long_desc = $4,
				start = $5,
				updated_at = now()
			where id = $1
		`, m.ID, m.Title, m.ShortDesc, m.LongDesc, pgsql.NullTime(&m.Start))
		if err != nil {
			return err
		}

		if imageURL != "" {
			err = dispatcher.Dispatch(ctx, &course.SetImage{ID: m.ID, Image: imageURL})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
