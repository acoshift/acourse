package course

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/acoshift/pgsql/pgctx"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/internal/pkg/image"
)

// Course model
type Course struct {
	ID          string
	Option      Option
	EnrollCount int64
	Title       string
	ShortDesc   string
	Desc        string
	Image       string
	Owner       struct {
		ID    string
		Name  string
		Image string
	}
	Start        pq.NullTime
	URL          string
	Type         int
	Price        float64
	Discount     float64
	EnrollDetail string
}

// Link returns id if url is invalid
func (x *Course) Link() string {
	if x.URL != "" {
		return x.URL
	}
	return x.ID
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
	if m.Title == "" {
		return "", fmt.Errorf("title required")
	}

	var imageURL string
	if m.Image != nil {
		err := image.Validate(m.Image)
		if err != nil {
			return "", err
		}

		img, err := m.Image.Open()
		if err != nil {
			return "", err
		}
		defer img.Close()

		imageURL, err = uploadCourseCoverImage(ctx, img)
		img.Close()
		if err != nil {
			return "", err
		}
	}

	var id string
	err := pgctx.RunInTx(ctx, func(ctx context.Context) error {
		// language=SQL
		err := pgctx.QueryRow(ctx, `
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
	if m.ID == "" {
		return fmt.Errorf("course id required")
	}
	if m.Title == "" {
		return fmt.Errorf("title required")
	}

	var imageURL string
	if m.Image != nil {
		err := image.Validate(m.Image)
		if err != nil {
			return err
		}

		img, err := m.Image.Open()
		if err != nil {
			return err
		}
		defer img.Close()

		imageURL, err = uploadCourseCoverImage(ctx, img)
		img.Close()
		if err != nil {
			return err
		}
	}

	err := pgctx.RunInTx(ctx, func(ctx context.Context) error {
		_, err := pgctx.Exec(ctx, `
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
	// language=SQL
	_, err := pgctx.Exec(ctx, `
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
	_, err := pgctx.Exec(ctx, `update courses set image = $2 where id = $1`, id, img)
	return err
}

// GetURL gets course url
func GetURL(ctx context.Context, id string) (string, error) {
	var r string
	err := pgctx.QueryRow(ctx,
		`select url from courses where id = $1`,
		id,
	).Scan(pgsql.NullString(&r))
	return r, err
}

// GetUserID gets course user id
func GetUserID(ctx context.Context, id string) (string, error) {
	var r string
	err := pgctx.QueryRow(ctx, `select user_id from courses where id = $1`, id).Scan(&r)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return r, err
}

// Get gets course from id
func Get(ctx context.Context, id string) (*Course, error) {
	var x Course
	// language=SQL
	err := pgctx.QueryRow(ctx, `
		select c. id, c.title, c.short_desc, c.long_desc, c.image,
		       c.start, c.url, c.type, c.price, c.discount, c.enroll_detail,
		       u.id, u.name, u.image,
		       opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from courses as c
			left join course_options as opt on opt.course_id = c.id
			left join users as u on u.id = c.user_id
		where c.id = $1
	`, id).Scan(
		&x.ID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image,
		&x.Start, pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Owner.ID, &x.Owner.Name, &x.Owner.Image,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

func GetIDByURL(ctx context.Context, url string) (courseID string, err error) {
	// language=SQL
	err = pgctx.QueryRow(ctx, `
		select id
		from courses
		where url = $1
	`, url).Scan(&courseID)
	if err == sql.ErrNoRows {
		err = ErrNotFound
	}
	return
}
