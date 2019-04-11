package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/acoshift/pgsql"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/internal/pkg/bus"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/email"
	"github.com/acoshift/acourse/internal/pkg/markdown"
	"github.com/acoshift/acourse/internal/pkg/model"
	"github.com/acoshift/acourse/internal/pkg/model/admin"
	"github.com/acoshift/acourse/internal/pkg/model/app"
	"github.com/acoshift/acourse/internal/pkg/model/course"
	"github.com/acoshift/acourse/internal/pkg/model/payment"
)

// Init inits admin service
func Init() {
	bus.Register(listUsers)
	bus.Register(countUsers)
	bus.Register(listCourses)
	bus.Register(countCourses)
	bus.Register(getPayment)
	bus.Register(listPayments)
	bus.Register(countPayments)
	bus.Register(acceptPayment)
	bus.Register(rejectPayment)
}

func listUsers(ctx context.Context, m *admin.ListUsers) error {
	rows, err := sqlctx.Query(ctx, `
		select
			id, name, username, email,
			image, created_at
		from users
		order by created_at desc
		limit $1 offset $2
	`, m.Limit, m.Offset)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var x admin.UserItem
		err = rows.Scan(
			&x.ID, &x.Name, &x.Username, pgsql.NullString(&x.Email),
			&x.Image, &x.CreatedAt,
		)
		if err != nil {
			return err
		}
		m.Result = append(m.Result, &x)
	}
	return rows.Err()
}

func countUsers(ctx context.Context, m *admin.CountUsers) error {
	return sqlctx.QueryRow(ctx,
		`select count(*) from users`,
	).Scan(&m.Result)
}

func listCourses(ctx context.Context, m *admin.ListCourses) error {
	rows, err := sqlctx.Query(ctx, `
		select
			c.id, c.title, c.image,
			c.url, c.type, c.price, c.discount,
			c.created_at, c.updated_at,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount,
			u.id, u.username, u.image
		from courses as c
			left join course_options as opt on opt.course_id = c.id
			left join users as u on u.id = c.user_id
		order by c.created_at desc
		limit $1 offset $2
	`, m.Limit, m.Offset)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var x admin.CourseItem
		err = rows.Scan(
			&x.ID, &x.Title, &x.Image,
			pgsql.NullString(&x.URL), &x.Type, &x.Price, &x.Discount,
			&x.CreatedAt, &x.UpdatedAt,
			&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
			&x.Owner.ID, &x.Owner.Username, &x.Owner.Image,
		)
		if err != nil {
			return err
		}
		m.Result = append(m.Result, &x)
	}
	return rows.Err()
}

func countCourses(ctx context.Context, m *admin.CountCourses) error {
	return sqlctx.QueryRow(ctx,
		`select count(*) from courses`,
	).Scan(&m.Result)
}

func getPayment(ctx context.Context, m *admin.GetPayment) error {
	x := &m.Result
	err := sqlctx.QueryRow(ctx, `
		select
			p.id,
			p.image, p.price, p.original_price, p.code,
			p.status, p.created_at, p.at,
			u.id, u.username, u.name, u.email, u.image,
			c.id, c.title, c.image, c.url
		from payments as p
			left join users as u on p.user_id = u.id
			left join courses as c on p.course_id = c.id
		where p.id = $1
	`, m.PaymentID).Scan(
		&x.ID,
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code,
		&x.Status, &x.CreatedAt, pgsql.NullTime(&x.At),
		&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
		&x.Course.ID, &x.Course.Title, &x.Course.Image, pgsql.NullString(&x.Course.URL),
	)
	if err == sql.ErrNoRows {
		return model.ErrNotFound
	}
	return err
}

func listPayments(ctx context.Context, m *admin.ListPayments) error {
	rows, err := sqlctx.Query(ctx, `
		select
			p.id,
			p.image, p.price, p.original_price, p.code,
			p.status, p.created_at, p.at,
			u.id, u.username, u.name, u.email, u.image,
			c.id, c.title, c.image, c.url
		from payments as p
			left join users as u on p.user_id = u.id
			left join courses as c on p.course_id = c.id
		where p.status = any($1)
		order by p.created_at desc
		limit $2 offset $3
	`, pq.Array(m.Status), m.Limit, m.Offset)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var x admin.Payment
		err = rows.Scan(
			&x.ID,
			&x.Image, &x.Price, &x.OriginalPrice, &x.Code,
			&x.Status, &x.CreatedAt, pgsql.NullTime(&x.At),
			&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
			&x.Course.ID, &x.Course.Title, &x.Course.Image, pgsql.NullString(&x.Course.URL),
		)
		if err != nil {
			return err
		}
		m.Result = append(m.Result, &x)
	}
	return rows.Err()
}

func countPayments(ctx context.Context, m *admin.CountPayments) error {
	return sqlctx.QueryRow(ctx, `
		select count(*)
		from payments
		where status = any($1)
	`, pq.Array(m.Status)).Scan(&m.Result)
}

func acceptPayment(ctx context.Context, m *admin.AcceptPayment) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		p := admin.GetPayment{PaymentID: m.ID}
		err := bus.Dispatch(ctx, &p)
		if err == model.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		err = bus.Dispatch(ctx, &payment.SetStatus{ID: p.Result.ID, Status: payment.Accepted})
		if err != nil {
			return err
		}

		return bus.Dispatch(ctx, &course.InsertEnroll{ID: p.Result.Course.ID, UserID: p.Result.User.ID})
	})
	if err != nil {
		return err
	}

	go func() {
		// re-fetch payment to get latest timestamp
		p := admin.GetPayment{PaymentID: m.ID}
		err := bus.Dispatch(ctx, &p)
		if err != nil {
			return
		}

		name := p.Result.User.Name
		if len(name) == 0 {
			name = p.Result.User.Username
		}
		body := markdown.Email(fmt.Sprintf(`สวัสดีครับคุณ %s,


อีเมล์ฉบับนี้ยืนยันว่าท่านได้รับการอนุมัติการชำระเงินสำหรับหลักสูตร "%s" เสร็จสิ้น ท่านสามารถทำการ login เข้าสู่ Website Acourse แล้วเข้าเรียนหลักสูตร "%s" ได้ทันที


รหัสการชำระเงิน: %s

ชื่อหลักสูตร: %s

จำนวนเงิน: %.2f บาท

เวลาที่ทำการชำระเงิน: %s

เวลาที่อนุมัติการชำระเงิน: %s

ชื่อผู้ชำระเงิน: %s

อีเมล์ผู้ชำระเงิน: %s

----------------------

ขอบคุณที่ร่วมเรียนกับเราครับ

ทีมงาน acourse.io

https://acourse.io
`,
			name,
			p.Result.Course.Title,
			p.Result.Course.Title,
			p.Result.ID,
			p.Result.Course.Title,
			p.Result.Price,
			p.Result.CreatedAt.In(m.Location).Format("02/01/2006 15:04:05"),
			p.Result.At.In(m.Location).Format("02/01/2006 15:04:05"),
			name,
			p.Result.User.Email,
		))

		title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", p.Result.Course.Title)
		email.Send(p.Result.User.Email, title, body)
	}()

	return nil
}

func rejectPayment(ctx context.Context, m *admin.RejectPayment) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		p := admin.GetPayment{PaymentID: m.ID}
		err := bus.Dispatch(ctx, &p)
		if err == model.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		return bus.Dispatch(ctx, &payment.SetStatus{ID: p.Result.ID, Status: payment.Rejected})
	})
	if err != nil {
		return err
	}

	go func() {
		p := admin.GetPayment{PaymentID: m.ID}
		err := bus.Dispatch(ctx, &p)
		if err != nil {
			return
		}
		body := markdown.Email(m.Message)
		title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", p.Result.Course.Title)
		email.Send(p.Result.User.Email, title, body)
	}()

	return nil
}
