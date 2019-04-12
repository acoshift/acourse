package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/acoshift/pgsql"
	"github.com/lib/pq"

	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/config"
	"github.com/acoshift/acourse/internal/pkg/context/sqlctx"
	"github.com/acoshift/acourse/internal/pkg/course"
	"github.com/acoshift/acourse/internal/pkg/email"
	"github.com/acoshift/acourse/internal/pkg/markdown"
	"github.com/acoshift/acourse/internal/pkg/payment"
)

// Payment type
type Payment struct {
	ID            string
	Image         string
	Price         float64
	OriginalPrice float64
	Code          string
	Status        int
	CreatedAt     time.Time
	At            time.Time
	User          struct {
		ID       string
		Username string
		Name     string
		Email    string
		Image    string
	}
	Course struct {
		ID    string
		Title string
		Image string
		URL   string
	}
}

// CourseLink returns course link
func (x *Payment) CourseLink() string {
	if x.Course.URL == "" {
		return x.Course.ID
	}
	return x.Course.URL
}

func GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	var x Payment
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
	`, paymentID).Scan(
		&x.ID,
		&x.Image, &x.Price, &x.OriginalPrice, &x.Code,
		&x.Status, &x.CreatedAt, pgsql.NullTime(&x.At),
		&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
		&x.Course.ID, &x.Course.Title, &x.Course.Image, pgsql.NullString(&x.Course.URL),
	)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func GetPayments(ctx context.Context, status []int, limit, offset int64) ([]*Payment, error) {
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
	`, pq.Array(status), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var xs []*Payment
	for rows.Next() {
		var x Payment
		err = rows.Scan(
			&x.ID,
			&x.Image, &x.Price, &x.OriginalPrice, &x.Code,
			&x.Status, &x.CreatedAt, pgsql.NullTime(&x.At),
			&x.User.ID, &x.User.Username, &x.User.Name, pgsql.NullString(&x.User.Email), &x.User.Image,
			&x.Course.ID, &x.Course.Title, &x.Course.Image, pgsql.NullString(&x.Course.URL),
		)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return xs, nil
}

func CountPayments(ctx context.Context, status []int) (cnt int64, err error) {
	err = sqlctx.QueryRow(ctx, `
		select count(*)
		from payments
		where status = any($1)
	`, pq.Array(status)).Scan(&cnt)
	return
}

func AcceptPayment(ctx context.Context, paymentID string) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		p, err := GetPayment(ctx, paymentID)
		if err == app.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		err = payment.SetStatus(ctx, p.ID, payment.Accepted)
		if err != nil {
			return err
		}

		return course.InsertEnroll(ctx, p.Course.ID, p.User.ID)
	})
	if err != nil {
		return err
	}

	go func() {
		// re-fetch payment to get latest timestamp
		p, err := GetPayment(ctx, paymentID)
		if err != nil {
			return
		}

		name := p.User.Name
		if len(name) == 0 {
			name = p.User.Username
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
			p.Course.Title,
			p.Course.Title,
			p.ID,
			p.Course.Title,
			p.Price,
			p.CreatedAt.In(config.Location()).Format("02/01/2006 15:04:05"),
			p.At.In(config.Location()).Format("02/01/2006 15:04:05"),
			name,
			p.User.Email,
		))

		title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", p.Course.Title)
		email.Send(p.User.Email, title, body)
	}()

	return nil
}

func RejectPayment(ctx context.Context, paymentID string, message string) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		p, err := GetPayment(ctx, paymentID)
		if err == app.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		return payment.SetStatus(ctx, p.ID, payment.Rejected)
	})
	if err != nil {
		return err
	}

	go func() {
		p, err := GetPayment(ctx, paymentID)
		if err != nil {
			return
		}
		body := markdown.Email(message)
		title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", p.Course.Title)
		email.Send(p.User.Email, title, body)
	}()

	return nil
}
