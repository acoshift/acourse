package admin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/acoshift/hime"
	"github.com/acoshift/paginate"

	"github.com/acoshift/acourse/context/sqlctx"
	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func (c *ctrl) rejectPayment(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	x, err := repository.GetPayment(ctx, id)
	if err != nil {
		return err
	}

	name := x.User.Name
	if len(name) == 0 {
		name = x.User.Username
	}
	message := fmt.Sprintf(`สวัสดีครับคุณ %s,


ตามที่ท่านได้ upload file เพื่อใช้ในการสมัครหลักสูตร "%s" เมื่อเวลา %s


ทางทีมงาน acourse.io ขอเรียนแจ้งให้ทราบว่าคำขอของคุณถูกปฏิเสธ โดยอาจจะเกิดจากสาเหตุใด สาเหตุหนึ่ง ตามรายละเอียดด้านล่าง


1. รูปภาพที่ upload ไม่ตรงกับสิ่งที่ระบุไว้ เช่น

  - สำหรับ Course free ไม่มีมัดจำ - รูปภาพต้องเป็นรูป screenshot จากการแชร์ link ของ course "https://acourse.io/course/%s" ไปยัง timeline facebook ของตนเองเท่านั้น
  - สำหรับ Course ประเภทอื่น ๆ ให้ลองอ่านรายละเอียดของรูปภาพที่จำเป็นต้องใช้ในการ upload ให้ครบถ้วนและปฏิบัติตามให้ถูกต้อง

1. จำนวนเงินที่ระบุไม่ตรงกับจำนวนเงินที่โอนจริง

  - ในกรณีที่ Course มีส่วนลด ให้ระบุยอดที่โอนเป็นตัวเลขที่ตรงกับยอดโอน เท่านั้น (ไม่ใช่ตัวเลขราคาเต็มของ Course)
  - ในกรณีที่จ่ายผ่าน 3rd party เช่น eventpop ให้ใส่ตามราคาบัตร ไม่รวมค่าบริการอื่น ๆ เช่นค่า fee ของ eventpop


ถ้าติดขัดหรือสงสัยตรงไหนเพิ่มเติม ท่านสามารถ reply email นี้เพื่อสอบถามเพิ่มเติมได้ครับ


ขอบคุณมากครับ

ทีมงาน acourse.io
`,
		name,
		x.Course.Title,
		x.CreatedAt.In(c.Location).Format("02/01/2006 15:04:05"),
		x.Course.Link(),
	)

	p := view.Page(ctx)
	p["Payment"] = x
	p["Message"] = message
	return ctx.View("admin.payment-reject", p)
}

func (c *ctrl) postRejectPayment(ctx *hime.Context) error {
	message := ctx.FormValue("message")
	id := ctx.FormValue("id")

	var x *entity.Payment
	err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
		var err error

		x, err = repository.GetPayment(ctx, id)
		if err != nil {
			return err
		}

		return repository.SetPaymentStatus(ctx, x.ID, entity.Rejected)
	})
	if err != nil {
		return err
	}

	if x.User.Email != "" {
		go func() {
			x, err := repository.GetPayment(ctx, id)
			if err != nil {
				return
			}
			body := view.Markdown(message)
			title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", x.Course.Title)
			c.EmailSender.Send(x.User.Email, title, body)
		}()
	}

	return ctx.RedirectTo("admin.payments.pending")
}

func (c *ctrl) postPendingPayment(ctx *hime.Context) error {
	action := ctx.FormValue("action")

	id := ctx.FormValue("id")
	if action == "accept" {
		var x *entity.Payment
		err := sqlctx.RunInTx(ctx, func(ctx context.Context) error {
			var err error
			x, err = repository.GetPayment(ctx, id)
			if err != nil {
				return err
			}

			err = repository.SetPaymentStatus(ctx, x.ID, entity.Accepted)
			if err != nil {
				return err
			}

			return repository.Enroll(ctx, x.UserID, x.CourseID)
		})
		if err != nil {
			return err
		}
		if x.User.Email != "" {
			go func() {
				// re-fetch payment to get latest timestamp
				x, err := repository.GetPayment(ctx, id)
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
					x.CreatedAt.In(c.Location).Format("02/01/2006 15:04:05"),
					x.At.Time.In(c.Location).Format("02/01/2006 15:04:05"),
					name,
					x.User.Email,
				))

				title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", x.Course.Title)
				c.EmailSender.Send(x.User.Email, title, body)
			}()
		}
	}
	return ctx.RedirectTo("admin.payments.pending")
}

func (c *ctrl) pendingPayments(ctx *hime.Context) error {
	cnt, err := repository.CountPaymentsByStatuses(ctx, []int{entity.Pending})
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	payments, err := repository.ListPaymentsByStatus(ctx, []int{entity.Pending}, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Navbar"] = "admin.payment.pending"
	p["Payments"] = payments
	p["Paginate"] = pn
	return ctx.View("admin.payments", p)
}

func (c *ctrl) historyPayments(ctx *hime.Context) error {
	cnt, err := repository.CountPaymentsByStatuses(ctx, []int{entity.Accepted, entity.Rejected})
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	payments, err := repository.ListPaymentsByStatus(ctx, []int{entity.Accepted, entity.Rejected, entity.Refunded}, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p["Navbar"] = "admin.payment.history"
	p["Payments"] = payments
	p["Paginate"] = pn
	return ctx.View("admin.payments", p)
}
