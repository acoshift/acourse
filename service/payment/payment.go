package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/moonrhythm/dispatcher"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/model/app"
	"github.com/acoshift/acourse/model/course"
	"github.com/acoshift/acourse/model/email"
	"github.com/acoshift/acourse/model/payment"
	"github.com/acoshift/acourse/view"
)

var (
	_loc *time.Location
)

// Init inits payment
func Init(loc *time.Location) {
	_loc = loc

	dispatcher.Register(setStatus)
	dispatcher.Register(hasPending)
	dispatcher.Register(accept)
	dispatcher.Register(reject)
}

func setStatus(ctx context.Context, m *payment.SetStatus) error {
	_, err := sqlctx.Exec(ctx, `
		update payments
		set
			status = $2,
			updated_at = now(),
			at = now()
		where id = $1
	`, m.ID, m.Status)
	return err
}

func hasPending(ctx context.Context, m *payment.HasPending) error {
	return sqlctx.QueryRow(ctx, `
		select exists (
			select 1
			from payments
			where user_id = $1 and course_id = $2 and status = $3
		)
	`, m.UserID, m.CourseID, entity.Pending).Scan(&m.Result)
}

func accept(ctx context.Context, m *payment.Accept) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		x, err := getPayment(ctx, m.ID)
		if err == entity.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		err = dispatcher.Dispatch(ctx, &payment.SetStatus{ID: x.ID, Status: entity.Accepted})
		if err != nil {
			return err
		}

		return dispatcher.Dispatch(ctx, &course.InsertEnroll{ID: x.Course.ID, UserID: x.User.ID})
	})
	if err != nil {
		return err
	}

	go func() {
		// re-fetch payment to get latest timestamp
		x, err := getPayment(ctx, m.ID)
		if err != nil {
			return
		}

		name := x.User.Name
		if len(name) == 0 {
			name = x.User.Username
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
			x.Course.Title,
			x.Course.Title,
			x.ID,
			x.Course.Title,
			x.Price,
			x.CreatedAt.In(_loc).Format("02/01/2006 15:04:05"),
			x.At.Time.In(_loc).Format("02/01/2006 15:04:05"),
			name,
			x.User.Email,
		))

		title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", x.Course.Title)
		dispatcher.Dispatch(context.Background(), &email.Send{
			To:      x.User.Email,
			Subject: title,
			Body:    body,
		})
	}()

	return nil
}

func reject(ctx context.Context, m *payment.Reject) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		x, err := getPayment(ctx, m.ID)
		if err == entity.ErrNotFound {
			return app.NewUIError("payment not found")
		}
		if err != nil {
			return err
		}

		return dispatcher.Dispatch(ctx, &payment.SetStatus{ID: x.ID, Status: entity.Rejected})
	})
	if err != nil {
		return err
	}

	go func() {
		x, err := getPayment(ctx, m.ID)
		if err != nil {
			return
		}
		body := view.MarkdownEmail(m.Message)
		title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", x.Course.Title)
		dispatcher.Dispatch(context.Background(), &email.Send{
			To:      x.User.Email,
			Subject: title,
			Body:    body,
		})
	}()

	return nil
}
