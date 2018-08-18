package repository

import (
	"context"
	"database/sql"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
)

// GetPayment gets payment from given id
func GetPayment(ctx context.Context, paymentID string) (*entity.Payment, error) {
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
		&x.ID,
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
