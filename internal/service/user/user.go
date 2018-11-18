package user

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/internal/context/sqlctx"
	"github.com/acoshift/acourse/internal/entity"
	"github.com/acoshift/acourse/internal/pkg/dispatcher"
	"github.com/acoshift/acourse/internal/pkg/model/user"
)

// Init inits user service
func Init() {
	dispatcher.Register(create)
	dispatcher.Register(update)
	dispatcher.Register(get)
	dispatcher.Register(isExists)
	dispatcher.Register(setImage)
	dispatcher.Register(updateProfile)
	dispatcher.Register(isEnroll)
	dispatcher.Register(enroll)
}

func create(ctx context.Context, m *user.Create) error {
	_, err := sqlctx.Exec(ctx, `
		insert into users
			(id, username, name, email, image)
		values
			($1, $2, $3, $4, $5)
	`, m.ID, m.Username, m.Name, pgsql.NullString(&m.Email), m.Image)
	if pgsql.IsUniqueViolation(err, "users_email_key") {
		return entity.ErrEmailNotAvailable
	}
	if pgsql.IsUniqueViolation(err, "users_username_key") {
		return entity.ErrUsernameNotAvailable
	}
	return err
}

func update(ctx context.Context, m *user.Update) error {
	_, err := sqlctx.Exec(ctx, `
		update users
		set
			username = $2,
			name = $3,
			about_me = $4,
			updated_at = now()
		where id = $1
	`, m.ID, m.Username, m.Name, m.AboutMe)
	return err
}

func get(ctx context.Context, m *user.Get) error {
	var x user.User
	err := sqlctx.QueryRow(ctx, `
		select
			u.id, u.name, u.username, coalesce(u.email, ''), u.about_me, u.image,
			coalesce(r.admin, false), coalesce(r.instructor, false)
		from users as u
			left join roles as r on u.id = r.user_id
		where u.id = $1
	`, m.ID).Scan(
		&x.ID, &x.Name, &x.Username, &x.Email, &x.AboutMe, &x.Image,
		&x.Role.Admin, &x.Role.Instructor,
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

func isExists(ctx context.Context, m *user.IsExists) error {
	return sqlctx.QueryRow(ctx, `
		select exists (
			select 1
			from users
			where id = $1
		)
	`, m.ID).Scan(&m.Result)
}

func setImage(ctx context.Context, m *user.SetImage) error {
	_, err := sqlctx.Exec(ctx, `
		update users
		set image = $2
		where id = $1
	`, m.ID, m.Image)
	return err
}

func isEnroll(ctx context.Context, m *user.IsEnroll) error {
	return sqlctx.QueryRow(ctx, `
		select exists (
			select 1
			from enrolls
			where user_id = $1 and course_id = $2
		)
	`, m.ID, m.CourseID).Scan(&m.Result)
}
