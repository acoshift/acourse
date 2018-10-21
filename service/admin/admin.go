package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/admin"
	"github.com/acoshift/acourse/model/app"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/model/email"
	"github.com/acoshift/acourse/model/payment"
	"github.com/acoshift/acourse/view"
	"github.com/acoshift/pgsql"
)

// Init inits admin service
func Init() {
	dispatcher.Register(acceptPayment)
	dispatcher.Register(rejectPayment)
	dispatcher.Register(getPayment)
}

func acceptPayment(ctx context.Context, m *admin.AcceptPayment) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		p := getPaymentModel{PaymentID: m.ID}
		err := dispatcher.Dispatch(ctx, &p)
		if err == entity.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		err = dispatcher.Dispatch(ctx, &payment.SetStatus{ID: p.Result.ID, Status: entity.Accepted})
		if err != nil {
			return err
		}

		return dispatcher.Dispatch(ctx, &course.InsertEnroll{ID: p.Result.Course.ID, UserID: p.Result.User.ID})
	})
	if err != nil {
		return err
	}

	go func() {
		// re-fetch payment to get latest timestamp
		p := getPaymentModel{PaymentID: m.ID}
		err := dispatcher.Dispatch(ctx, &p)
		if err != nil {
			return
		}

		name := p.Result.User.Name
		if len(name) == 0 {
			name = p.Result.User.Username
		}
		body := view.MarkdownEmail(fmt.Sprintf(`สวัสดีครับคุณ %s,


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
		dispatcher.Dispatch(context.Background(), &email.Send{
			To:      p.Result.User.Email,
			Subject: title,
			Body:    body,
		})
	}()

	return nil
}

func rejectPayment(ctx context.Context, m *admin.RejectPayment) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		p := getPaymentModel{PaymentID: m.ID}
		err := dispatcher.Dispatch(ctx, &p)
		if err == entity.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		return dispatcher.Dispatch(ctx, &payment.SetStatus{ID: p.Result.ID, Status: entity.Rejected})
	})
	if err != nil {
		return err
	}

	go func() {
		p := getPaymentModel{PaymentID: m.ID}
		err := dispatcher.Dispatch(ctx, &p)
		if err != nil {
			return
		}
		body := view.MarkdownEmail(m.Message)
		title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", p.Result.Course.Title)
		dispatcher.Dispatch(context.Background(), &email.Send{
			To:      p.Result.User.Email,
			Subject: title,
			Body:    body,
		})
	}()

	return nil
}

type getPaymentModel struct {
	PaymentID string

	Result struct {
		ID        string
		Price     float64
		Status    int
		CreatedAt time.Time
		At        time.Time

		User struct {
			ID       string
			Username string
			Name     string
			Email    string
		}
		Course struct {
			ID    string
			Title string
		}
	}
}

func getPayment(ctx context.Context, m *getPaymentModel) error {
	r := &m.Result
	err := sqlctx.QueryRow(ctx, `
		select
			p.id,
			p.price,
			p.status, p.created_at, p.at,
			u.id, u.username, u.name, u.email,
			c.id, c.title
		from payments as p
			left join users as u on p.user_id = u.id
			left join courses as c on p.course_id = c.id
		where p.id = $1
	`, m.PaymentID).Scan(
		&r.ID,
		&r.Price,
		&r.Status, &r.CreatedAt, pgsql.NullTime(&r.At),
		&r.User.ID, &r.User.Username, &r.User.Name, pgsql.NullString(&r.User.Email),
		&r.Course.ID, &r.Course.Title,
	)
	if err == sql.ErrNoRows {
		return entity.ErrNotFound
	}
	return err
}
