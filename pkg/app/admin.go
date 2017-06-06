package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/view"
)

func adminUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	cnt, err := model.CountUsers()
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

	users, err := model.ListUsers(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.AdminUsers(w, r, users, int(page), int(totalPage))
}

func getAdminCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := model.ListCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.AdminCourses(w, r, courses, 1, 1)
}

func adminPayments(w http.ResponseWriter, r *http.Request, paymentsGetter func(int64, int64) ([]*model.Payment, error), paymentsCounter func() (int64, error)) {
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if page <= 0 {
		page = 1
	}
	limit := int64(30)

	cnt, err := paymentsCounter()
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

	payments, err := paymentsGetter(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.AdminPayments(w, r, payments, int(page), int(totalPage))
}

func postAdminPendingPayment(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("Action")

	id, err := strconv.ParseInt(r.FormValue("ID"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if action == "accept" {
		x, err := model.GetPayment(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = x.Accept()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if x.User.Email.Valid {
			go func() {
				// re-fetch payment to get latest timestamp
				x, err := model.GetPayment(id)
				if err != nil {
					return
				}

				name := x.User.Name
				if len(name) == 0 {
					name = x.User.Username
				}
				body := fmt.Sprintf(`สวัสดีครับคุณ %s,<br>
<br>
อีเมล์ฉบับนี้ยืนยันว่าท่านได้รับการอนุมัติการชำระเงินสำหรับหลักสูตร "%s" เสร็จสิ้น ท่านสามารถทำการ login เข้าสู่ Website Acourse แล้วเข้าเรียนหลักสูตร "%s" ได้ทันที<br>
<br>
รหัสการชำระเงิน: %d<br>
ชื่อหลักสูตร: %s<br>
จำนวนเงิน: %.2f บาท<br>
เวลาที่ทำการชำระเงิน: %s<br>
เวลาที่อนุมัติการชำระเงิน: %s<br>
ชื่อผู้ชำระเงิน: %s<br>
อีเมล์ผู้ชำระเงิน: %s<br>
----------------------<br>
ขอบคุณที่ร่วมเรียนกับเราครับ<br>
Krissada Chalermsook (Oak)<br>
Founder/CEO Acourse.io<br>
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
				)

				title := fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", x.Course.Title)
				SendEmail(x.User.Email.String, title, body)
			}()
		}
	} else if action == "reject" {
		x, err := model.GetPayment(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = x.Reject()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if x.User.Email.Valid {
			go func() {
				x, err := model.GetPayment(id)
				if err != nil {
					return
				}
				name := x.User.Name
				if len(name) == 0 {
					name = x.User.Username
				}

				body := fmt.Sprintf(`สวัสดีครับคุณ %s,<br>
<br>
ตามที่ท่านได้ upload file เพื่อใช้ในการสมัครหลักสูตร "%s" เมื่อเวลา %s<br>
<br>
ทางทีมงาน acourse.io ขอเรียนแจ้งให้ทราบว่าคำขอของคุณถูกปฏิเสธ โดยอาจจะเกิดจากสาเหตุใด สาเหตุหนึ่ง ตามรายละเอียดด้านล่าง<br>
<br>
1.รูปภาพที่ upload ไม่ตรงกับสิ่งที่ระบุไว้ เช่น<br>
- สำหรับ Course free ไม่มีมัดจำ - รูปภาพต้องเป็นรูป screenshot จากการแชร์ link ของ course "https://acourse.io/course/%s" ไปยัง timeline facebook ของตนเองเท่านั้น<br>
- สำหรับ Course ประเภทอื่นๆ ให้ลองอ่านรายละเอียดของรูปภาพที่จำเป็นต้องใช้ในการ upload ให้ครบถ้วนและปฏิบัติตามให้ถูกต้อง<br>
<br>
2. จำนวนเงินที่ระบุไม่ตรงกับจำนวนเงินที่โอนจริง<br>
- ในกรณีที่ Course มีส่วนลด ให้ระบุยอดที่โอนเป็นตัวเลขที่ตรงกับยอดโอน เท่านั้น (ไม่ใช่ตัวเลขราคาเต็มของ Course)<br>
- ในกรณีที่จ่ายผ่าน 3rd party เช่น eventpop ให้ใส่ตามราคาบัตร ไม่รวมค่าบริการอื่น ๆ เช่นค่า fee ของ eventpop<br>
<br>
ถ้าติดขัดหรือสงสัยตรงไหนเพิ่มเติม ท่านสามารถ reply email นี้เพื่อสอบถามเพิ่มเติมได้ครับ<br>
<br>
ขอบคุณมากครับ<br>
ทีมงาน acourse.io
`,
					name,
					x.Course.Title,
					x.CreatedAt.In(loc).Format("02/01/2006 15:04:05"),
					x.Course.Link(),
				)

				title := fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", x.Course.Title)
				SendEmail(x.User.Email.String, title, body)
			}()
		}
	}
	http.Redirect(w, r, "/admin/payments/pending", http.StatusSeeOther)
}

func adminPendingPayments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		adminPayments(w, r, model.ListPendingPayments, model.CountPendingPayments)
	} else if r.Method == http.MethodPost {
		postAdminPendingPayment(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func adminHistoryPayments(w http.ResponseWriter, r *http.Request) {
	adminPayments(w, r, model.ListHistoryPayments, model.CountHistoryPayments)
}
