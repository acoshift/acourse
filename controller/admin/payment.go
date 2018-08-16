package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/acoshift/hime"
	"github.com/acoshift/paginate"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/service"
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
	id := ctx.FormValue("id")
	message := ctx.PostFormValue("message")

	err := c.Service.RejectPayment(ctx, id, message)
	if service.IsUIError(err) {
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	return ctx.RedirectTo("admin.payments.pending")
}

func (c *ctrl) postPendingPayment(ctx *hime.Context) error {
	action := ctx.FormValue("action")

	id := ctx.PostFormValue("id")
	if action == "accept" {
		err := c.Service.AcceptPayment(ctx, id)
		if service.IsUIError(err) {
			return ctx.Status(http.StatusBadRequest).String(err.Error())
		}
		if err != nil {
			return err
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
