package controller

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/acoshift/acourse/pkg/app"
)

func (c *ctrl) AdminUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	cnt, err := c.repo.CountUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	for offset > cnt {
		page--
		offset = (page - 1) * limit
	}
	totalPage := cnt / limit

	users, err := c.repo.ListUsers(ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.view.AdminUsers(w, r, users, int(page), int(totalPage))
}

func (c *ctrl) AdminCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	cnt, err := c.repo.CountCourses(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	for offset > cnt {
		page--
		offset = (page - 1) * limit
	}
	totalPage := int64(math.Ceil(float64(cnt) / float64(limit)))

	courses, err := c.repo.ListCourses(ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.view.AdminCourses(w, r, courses, int(page), int(totalPage))
}

func (c *ctrl) adminPayments(w http.ResponseWriter, r *http.Request, paymentsGetter func(context.Context, int64, int64) ([]*app.Payment, error), paymentsCounter func(context.Context) (int64, error)) {
	ctx := r.Context()
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	cnt, err := paymentsCounter(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	for offset > cnt {
		page--
		offset = (page - 1) * limit
	}
	totalPage := cnt / limit

	payments, err := paymentsGetter(ctx, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.view.AdminPayments(w, r, payments, int(page), int(totalPage))
}

func (c *ctrl) AdminRejectPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postAdminRejectPayment(w, r)
		return
	}

	ctx := r.Context()
	id := r.FormValue("id")
	x, err := c.repo.GetPayment(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		x.CreatedAt.In(c.loc).Format("02/01/2006 15:04:05"),
		x.Course.Link(),
	)

	c.view.AdminPaymentReject(w, r, x, message)
}

func (c *ctrl) postAdminRejectPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	message := r.FormValue("Message")
	id := r.FormValue("ID")

	x, err := c.repo.GetPayment(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = c.repo.RejectPayment(ctx, x)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if x.User.Email.Valid {
		go func() {
			x, err := c.repo.GetPayment(ctx, id)
			if err != nil {
				return
			}
			body := markdown(message)
			title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", x.Course.Title)
			c.sendEmail(x.User.Email.String, title, body)
		}()
	}

	http.Redirect(w, r, "/admin/payments/pending", http.StatusSeeOther)
}

func (c *ctrl) postAdminPendingPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	action := r.FormValue("Action")

	id := r.FormValue("ID")
	if action == "accept" {
		txctx, tx, err := app.WithTransaction(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		x, err := c.repo.GetPayment(txctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = c.repo.AcceptPayment(txctx, x)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tx.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if x.User.Email.Valid {
			go func() {
				// re-fetch payment to get latest timestamp
				x, err := c.repo.GetPayment(ctx, id)
				if err != nil {
					return
				}

				name := x.User.Name
				if len(name) == 0 {
					name = x.User.Username
				}
				body := markdown(fmt.Sprintf(`สวัสดีครับคุณ %s,


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
					x.CreatedAt.In(c.loc).Format("02/01/2006 15:04:05"),
					x.At.Time.In(c.loc).Format("02/01/2006 15:04:05"),
					name,
					x.User.Email.String,
				))

				title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", x.Course.Title)
				c.sendEmail(x.User.Email.String, title, body)
			}()
		}
	}
	http.Redirect(w, r, "/admin/payments/pending", http.StatusSeeOther)
}

func (c *ctrl) AdminPendingPayments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postAdminPendingPayment(w, r)
		return
	}
	c.adminPayments(w, r, c.repo.ListPendingPayments, c.repo.CountPendingPayments)
}

func (c *ctrl) AdminHistoryPayments(w http.ResponseWriter, r *http.Request) {
	c.adminPayments(w, r, c.repo.ListHistoryPayments, c.repo.CountHistoryPayments)
}
