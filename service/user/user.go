package user

import (
	"context"

	"github.com/acoshift/pgsql"
	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/user"
)

// Init inits user service
func Init() {
	dispatcher.Register(create)
	dispatcher.Register(update)
	dispatcher.Register(isExists)
	dispatcher.Register(setImage)
	dispatcher.Register(updateProfile)
}

func create(ctx context.Context, m *user.Create) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
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
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
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

func isExists(ctx context.Context, m *user.IsExists) error {
	q := sqlctx.GetQueryer(ctx)

	return q.QueryRow(`
		select exists (
			select 1
			from users
			where id = $1
		)
	`, m.ID).Scan(&m.Result)
}

func setImage(ctx context.Context, m *user.SetImage) error {
	q := sqlctx.GetQueryer(ctx)

	_, err := q.Exec(`
		update users
		set image = $2
		where id = $1
	`, m.ID, m.Image)
	return err
}
