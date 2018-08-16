package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/go-redis/redis"

	"github.com/acoshift/acourse/context/redisctx"
	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/service"
)

// NewService creates new service repository
func NewService() service.Repository {
	return &svcRepo{}
}

type svcRepo struct {
}

func (svcRepo) StoreMagicLink(ctx context.Context, linkID string, userID string) error {
	c := redisctx.GetClient(ctx)
	prefix := redisctx.GetPrefix(ctx)

	err := c.Set(prefix+"magic:"+linkID, userID, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (svcRepo) FindMagicLink(ctx context.Context, linkID string) (string, error) {
	c := redisctx.GetClient(ctx)
	prefix := redisctx.GetPrefix(ctx)

	key := prefix + "magic:" + linkID
	userID, err := c.Get(key).Result()
	if err == redis.Nil {
		return "", entity.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	c.Del(key)
	return userID, nil
}

func (svcRepo) CanAcquireMagicLink(ctx context.Context, email string) (bool, error) {
	c := redisctx.GetClient(ctx)
	prefix := redisctx.GetPrefix(ctx)

	key := prefix + "magic-rate:" + email
	current, err := c.Incr(key).Result()
	if err != nil {
		return false, err
	}
	if current > 1 {
		return false, nil
	}
	err = c.Expire(key, 5*time.Minute).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (svcRepo) GetEmailSignInUserByEmail(ctx context.Context, email string) (*entity.EmailSignInUser, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.EmailSignInUser
	err := q.QueryRow(`
		select
			id, name, email
		from users
		where email = $1
	`, email).Scan(
		&x.ID, &x.Name, pgsql.NullString(&x.Email),
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (svcRepo) RegisterUser(ctx context.Context, x *service.RegisterUser) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into users
			(id, username, name, email, image)
		values
			($1, $2, $3, $4, $5)
	`, x.ID, x.Username, x.Name, pgsql.NullString(&x.Email), x.Image)
	if pgsql.IsUniqueViolation(err, "users_email_key") {
		return entity.ErrEmailNotAvailable
	}
	if pgsql.IsUniqueViolation(err, "users_username_key") {
		return entity.ErrUsernameNotAvailable
	}
	return err
}

func (svcRepo) UpdateUser(ctx context.Context, x *service.UpdateUser) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update users
		set
			username = $2,
			name = $3,
			about_me = $4,
			updated_at = now()
		where id = $1
	`, x.ID, x.Username, x.Name, x.AboutMe)
	return err
}

func (svcRepo) SetUserImage(ctx context.Context, userID string, image string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update users
		set image = $2
		where id = $1
	`, userID, image)
	return err
}

func (svcRepo) IsUserExists(ctx context.Context, userID string) (exists bool, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select exists (
			select 1
			from users
			where id = $1
		)
	`, userID).Scan(&exists)
	return
}

func (svcRepo) RegisterCourse(ctx context.Context, x *service.RegisterCourse) (courseID string, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		insert into courses
			(user_id, title, short_desc, long_desc, image, start)
		values
			($1, $2, $3, $4, $5, $6)
		returning id
	`, x.UserID, x.Title, x.ShortDesc, x.LongDesc, x.Image, pgsql.NullTime(&x.Start)).Scan(&courseID)
	return
}

func (svcRepo) GetCourse(ctx context.Context, courseID string) (*entity.Course, error) {
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
		&x.Start, &x.URL, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
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

func (svcRepo) UpdateCourse(ctx context.Context, x *service.UpdateCourseModel) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update courses
		set
			title = $2,
			short_desc = $3,
			long_desc = $4,
			start = $5,
			updated_at = now()
		where id = $1
	`, x.ID, x.Title, x.ShortDesc, x.LongDesc, pgsql.NullTime(&x.Start))
	return err
}

func (svcRepo) SetCourseImage(ctx context.Context, courseID string, image string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`update courses set image = $2 where id = $1`, courseID, image)
	return err
}

func (svcRepo) SetCourseOption(ctx context.Context, courseID string, x *entity.CourseOption) error {
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
	`, courseID, x.Public, x.Enroll, x.Attend, x.Assignment, x.Discount)
	return err
}

func (svcRepo) RegisterPayment(ctx context.Context, x *service.RegisterPayment) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into payments
			(user_id, course_id, image, price, original_price, code, status)
		values
			($1, $2, $3, $4, $5, $6, $7)
		returning id
	`, x.UserID, x.CourseID, x.Image, x.Price, x.OriginalPrice, x.Code, x.Status)
	if err != nil {
		return err
	}
	return nil
}

func (svcRepo) GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error) {
	q := sqlctx.GetQueryer(ctx)

	var x entity.Payment
	err := q.QueryRow(`
		select
			payments.id,
			payments.image,
			payments.price,
			payments.original_price,
			payments.code,
			payments.status,
			payments.created_at,
			payments.updated_at,
			payments.at,
			users.id,
			users.username,
			users.name,
			users.email,
			users.image,
			courses.id,
			courses.title,
			courses.image,
			courses.url
		from payments
			left join users on payments.user_id = users.id
			left join courses on payments.course_id = courses.id
		where payments.id = $1
	`, paymentID).Scan(
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code, &x.Status, &x.CreatedAt, &x.UpdatedAt, &x.At,
		&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
		&x.Course.ID, &x.Course.Title, &x.Course.Image, &x.Course.URL,
	)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	x.UserID = x.User.ID
	x.CourseID = x.Course.ID
	return &x, nil
}

func (svcRepo) SetPaymentStatus(ctx context.Context, paymentID string, status int) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, paymentID, status)
	return err
}

func (svcRepo) HasPendingPayment(ctx context.Context, userID string, courseID string) (exists bool, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select exists (
			select 1
			from payments
			where user_id = $1 and course_id = $2 and status = $3
		)
	`, userID, courseID, entity.Pending).Scan(&exists)
	return
}

func (svcRepo) RegisterEnroll(ctx context.Context, userID string, courseID string) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		insert into enrolls
			(user_id, course_id)
		values
			($1, $2)
	`, userID, courseID)
	return err
}

func (svcRepo) IsEnrolled(ctx context.Context, userID string, courseID string) (enrolled bool, err error) {
	q := sqlctx.GetQueryer(ctx)

	err = q.QueryRow(`
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, userID, courseID).Scan(&enrolled)
	return
}
