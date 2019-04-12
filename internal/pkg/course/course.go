package course

import (
	"context"
	"database/sql"
	"mime/multipart"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/image"
)

// Course model
type Course struct {
	ID            string
	Option        Option
	EnrollCount   int64
	Title         string
	ShortDesc     string
	Desc          string
	Image         string
	UserID        string
	Start         pq.NullTime
	URL           string
	Type          int
	Price         float64
	Discount      float64
	Contents      []*Content
	EnrollDetail  string
	AssignmentIDs []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Link returns id if url is invalid
func (x *Course) Link() string {
	if x.URL == "" {
		return x.ID
	}
	return x.URL
}

// Option type
type Option struct {
	Public     bool
	Enroll     bool
	Attend     bool
	Assignment bool
	Discount   bool
}

// Course type values
const (
	_ = iota
	Live
	Video
	EBook
)

// Video type values
const (
	_ = iota
	Youtube
)

type CreateArgs struct {
	UserID    string
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time
}

// Create creates new course
func Create(ctx context.Context, m *CreateArgs) (string, error) {
	// TODO: validate user role

	if m.Title == "" {
		return "", app.NewUIError("title required")
	}

	var imageURL string
	if m.Image != nil {
		err := image.Validate(m.Image)
		if err != nil {
			return "", err
		}

		image, err := m.Image.Open()
		if err != nil {
			return "", err
		}
		defer image.Close()

		imageURL, err = uploadCourseCoverImage(ctx, image)
		image.Close()
		if err != nil {
			return "", app.NewUIError(err.Error())
		}
	}

	var id string
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		err := sqlctx.QueryRow(ctx, `
			insert into courses
				(user_id, title, short_desc, long_desc, image, start)
			values
				($1, $2, $3, $4, $5, $6)
			returning id
		`, m.UserID, m.Title, m.ShortDesc, m.LongDesc, imageURL, pgsql.NullTime(&m.Start)).Scan(&id)
		if err != nil {
			return err
		}

		return SetOption(ctx, id, Option{})
	})
	return id, err
}

type UpdateArgs struct {
	ID        string
	Title     string
	ShortDesc string
	LongDesc  string
	Image     *multipart.FileHeader
	Start     time.Time
}

// Update updates course
func Update(ctx context.Context, m *UpdateArgs) error {
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
			err = SetImage(ctx, m.ID, imageURL)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// SetOption sets course option
func SetOption(ctx context.Context, id string, opt Option) error {
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
	`,
		id,
		opt.Public, opt.Enroll, opt.Attend, opt.Assignment, opt.Discount,
	)
	return err
}

// SetImage sets course image
func SetImage(ctx context.Context, id string, img string) error {
	_, err := sqlctx.Exec(ctx, `update courses set image = $2 where id = $1`, id, img)
	return err
}

// GetURL gets course url
func GetURL(ctx context.Context, id string) (string, error) {
	var r string
	err := sqlctx.QueryRow(ctx,
		`select url from courses where id = $1`,
		id,
	).Scan(pgsql.NullString(&r))
	return r, err
}

// GetUserID gets course user id
func GetUserID(ctx context.Context, id string) (string, error) {
	var r string
	err := sqlctx.QueryRow(ctx, `select user_id from courses where id = $1`, id).Scan(&r)
	if err == sql.ErrNoRows {
		return "", app.ErrNotFound
	}
	return r, err
}

// Get gets course from id
func Get(ctx context.Context, id string) (*Course, error) {
	var x Course
	err := sqlctx.QueryRow(ctx, `
		select
			id, user_id, title, short_desc, long_desc, image,
			start, url, type, price, courses.discount, enroll_detail,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from courses
			left join course_options as opt on opt.course_id = courses.id
		where id = $1
	`, id).Scan(
		&x.ID, &x.UserID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image,
		&x.Start, pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

// InsertEnroll inserts enroll
func InsertEnroll(ctx context.Context, id, userID string) error {
	_, err := sqlctx.Exec(ctx, `
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, userID, id)
	return err
}
