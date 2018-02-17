package app

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/acoshift/pgsql"

	"github.com/acoshift/acourse/entity"
	"github.com/acoshift/acourse/repository"
	"github.com/acoshift/acourse/view"
)

func adminUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	cnt, err := repository.CountUsers(db)
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

	users, err := repository.ListUsers(db, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.AdminUsers(w, r, users, int(page), int(totalPage))
}

func adminCourses(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	cnt, err := repository.CountCourses(db)
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

	courses, err := repository.ListCourses(db, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.AdminCourses(w, r, courses, int(page), int(totalPage))
}

func adminPayments(w http.ResponseWriter, r *http.Request, history bool) {
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	var err error
	var cnt int64
	if history {
		cnt, err = repository.CountHistoryPayments(db)
	} else {
		cnt, err = repository.CountPendingPayments(db)
	}
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

	var payments []*entity.Payment
	if history {
		payments, err = repository.ListHistoryPayments(db, limit, offset)
	} else {
		payments, err = repository.ListPendingPayments(db, limit, offset)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.AdminPayments(w, r, payments, int(page), int(totalPage))
}

func adminRejectPayment(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	x, err := repository.GetPayment(db, id)
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
		x.CreatedAt.In(loc).Format("02/01/2006 15:04:05"),
		x.Course.Link(),
	)

	view.AdminPaymentReject(w, r, x, message)
}

func postAdminRejectPayment(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("Message")
	id := r.FormValue("ID")

	var x *entity.Payment
	err := pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
		var err error
		x, err = repository.GetPayment(tx, id)
		if err != nil {
			return err
		}
		return repository.RejectPayment(tx, x)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if x.User.Email.Valid {
		go func() {
			x, err := repository.GetPayment(db, id)
			if err != nil {
				return
			}
			body := markdown(message)
			title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", x.Course.Title)
			sendEmail(x.User.Email.String, title, body)
		}()
	}

	http.Redirect(w, r, "/admin/payments/pending", http.StatusSeeOther)
}

func postAdminPendingPayment(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("Action")

	id := r.FormValue("ID")
	if action == "accept" {
		var x *entity.Payment
		err := pgsql.RunInTx(db, nil, func(tx *sql.Tx) error {
			var err error
			x, err = repository.GetPayment(tx, id)
			if err != nil {
				return err
			}
			return repository.AcceptPayment(tx, x)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if x.User.Email.Valid {
			go func() {
				// re-fetch payment to get latest timestamp
				x, err := repository.GetPayment(db, id)
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
					x.CreatedAt.In(loc).Format("02/01/2006 15:04:05"),
					x.At.Time.In(loc).Format("02/01/2006 15:04:05"),
					name,
					x.User.Email.String,
				))

				title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", x.Course.Title)
				sendEmail(x.User.Email.String, title, body)
			}()
		}
	}
	http.Redirect(w, r, "/admin/payments/pending", http.StatusSeeOther)
}

func adminPendingPayments(w http.ResponseWriter, r *http.Request) {
	adminPayments(w, r, false)
}

func adminHistoryPayments(w http.ResponseWriter, r *http.Request) {
	adminPayments(w, r, true)
}
