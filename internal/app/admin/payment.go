package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/moonrhythm/hime"

	"github.com/acoshift/acourse/internal/app/view"
	"github.com/acoshift/acourse/internal/pkg/admin"
	"github.com/acoshift/acourse/internal/pkg/app"
	"github.com/acoshift/acourse/internal/pkg/config"
	"github.com/acoshift/acourse/internal/pkg/payment"
)

func getRejectPayment(ctx *hime.Context) error {
	id := ctx.FormValue("id")

	x, err := admin.GetPayment(ctx, id)
	if err == app.ErrNotFound {
		return ctx.RedirectTo("admin.payments.pending")
	}
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
		x.CreatedAt.In(config.Location()).Format("02/01/2006 15:04:05"),
		x.CourseLink(),
	)

	p := view.Page(ctx)
	p.Data["Payment"] = x
	p.Data["Message"] = message
	return ctx.View("admin.payment-reject", p)
}

func postRejectPayment(ctx *hime.Context) error {
	id := ctx.FormValue("id")
	message := ctx.PostFormValue("message")

	err := admin.RejectPayment(ctx, id, message)
	if app.IsUIError(err) {
		return ctx.Status(http.StatusBadRequest).String(err.Error())
	}
	if err != nil {
		return err
	}

	return ctx.RedirectTo("admin.payments.pending")
}

func postPendingPayment(ctx *hime.Context) error {
	action := ctx.FormValue("action")

	id := ctx.PostFormValue("id")
	if action == "accept" {
		err := admin.AcceptPayment(ctx, id)
		if app.IsUIError(err) {
			return ctx.Status(http.StatusBadRequest).String(err.Error())
		}
		if err != nil {
			return err
		}
	}

	return ctx.RedirectTo("admin.payments.pending")
}

func getPendingPayments(ctx *hime.Context) error {
	cnt, err := admin.CountPayments(ctx, []int{payment.Pending})
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	list, err := admin.GetPayments(ctx, []int{payment.Pending}, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.payment.pending"
	p.Data["Payments"] = list
	p.Data["Paginate"] = pn
	return ctx.View("admin.payments", p)
}

func getHistoryPayments(ctx *hime.Context) error {
	cnt, err := admin.CountPayments(ctx, []int{payment.Accepted, payment.Rejected})
	if err != nil {
		return err
	}

	pg, _ := strconv.ParseInt(ctx.FormValue("page"), 10, 64)
	pn := paginate.New(pg, 30, cnt)

	list, err := admin.GetPayments(ctx, []int{payment.Accepted, payment.Rejected, payment.Refunded}, pn.Limit(), pn.Offset())
	if err != nil {
		return err
	}

	p := view.Page(ctx)
	p.Data["Navbar"] = "admin.payment.history"
	p.Data["Payments"] = list
	p.Data["Paginate"] = pn
	return ctx.View("admin.payments", p)
}
