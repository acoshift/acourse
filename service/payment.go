package service

import (
	"context"
	"fmt"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func (s *svc) AcceptPayment(ctx context.Context, paymentID string) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		x, err := repository.GetPayment(ctx, paymentID)
		if err == entity.ErrNotFound {
			return newUIError("payment not found")
		}
		if err != nil {
			return err
		}

		err = repository.SetPaymentStatus(ctx, x.ID, entity.Accepted)
		if err != nil {
			return err
		}

		return repository.RegisterEnroll(ctx, x.UserID, x.CourseID)
	})
	if err != nil {
		return err
	}

	go func() {
		// re-fetch payment to get latest timestamp
		x, err := repository.GetPayment(ctx, paymentID)
		if err != nil {
			return
		}

		name := x.User.Name
		if len(name) == 0 {
			name = x.User.Username
		}
		body := view.Markdown(fmt.Sprintf(`สวัสดีครับคุณ %s,


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
			x.CreatedAt.In(s.Location).Format("02/01/2006 15:04:05"),
			x.At.Time.In(s.Location).Format("02/01/2006 15:04:05"),
			name,
			x.User.Email,
		))

		title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", x.Course.Title)
		s.EmailSender.Send(x.User.Email, title, body)
	}()

	return nil
}

func (s *svc) RejectPayment(ctx context.Context, paymentID string, msg string) error {
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		x, err := repository.GetPayment(ctx, paymentID)
		if err == entity.ErrNotFound {
			return newUIError("payment not found")
		}
		if err != nil {
			return err
		}

		return repository.SetPaymentStatus(ctx, x.ID, entity.Rejected)
	})
	if err != nil {
		return err
	}

	go func() {
		x, err := repository.GetPayment(ctx, paymentID)
		if err != nil {
			return
		}
		body := view.Markdown(msg)
		title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", x.Course.Title)
		s.EmailSender.Send(x.User.Email, title, body)
	}()

	return nil
}
