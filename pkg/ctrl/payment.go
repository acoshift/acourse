package ctrl

import (
	"fmt"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/e"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/acourse/pkg/view"
	"github.com/acoshift/go-firebase-admin"
)

// PaymentController implements PaymentController interface
type PaymentController struct {
	db   *store.DB
	auth *admin.FirebaseAuth
}

// NewPaymentController creates new controller
func NewPaymentController(db *store.DB, auth *admin.FirebaseAuth) *PaymentController {
	return &PaymentController{db, auth}
}

// List runs list action
func (c *PaymentController) List(ctx *app.PaymentListContext) (interface{}, error) {
	role, err := c.db.RoleGet(ctx.CurrentUserID)
	if err != nil {
		return nil, err
	}

	// only admin can access
	if !role.Admin {
		return nil, e.ErrForbidden
	}

	var xs []*model.Payment

	if ctx.History {
		xs, err = c.db.PaymentList()
	} else {
		xs, err = c.db.PaymentList(store.PaymentListOptionStatus(model.PaymentStatusWaiting))
	}
	if err != nil {
		return nil, err
	}

	res := make(view.PaymentCollection, len(xs))
	for i, x := range xs {
		user, err := c.db.UserMustGet(x.UserID)
		if err != nil {
			return nil, err
		}
		course, err := c.db.CourseGet(x.CourseID)
		if err != nil {
			return nil, err
		}
		res[i] = ToPaymentView(x, ToUserTinyView(user), ToCourseMiniView(course))
	}

	return res, nil
}

// Approve runs approve action
func (c *PaymentController) Approve(ctx *app.PaymentApproveContext) error {
	role, err := c.db.RoleGet(ctx.CurrentUserID)
	if err != nil {
		return err
	}
	if !role.Admin {
		return e.ErrForbidden
	}

	payment, err := c.db.PaymentGet(ctx.PaymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return e.ErrNotFound
	}
	payment.Approve()

	// Add user to enroll
	enroll := &model.Enroll{
		UserID:   payment.UserID,
		CourseID: payment.CourseID,
	}
	err = c.db.EnrollSave(enroll)
	if err != nil {
		return err
	}

	err = c.db.PaymentSave(payment)
	if err != nil {
		return err
	}

	go c.approved(payment)

	return nil
}

func (c *PaymentController) approved(payment *model.Payment) {
	course, err := c.db.CourseGet(payment.CourseID)
	if err != nil {
		log.Println(err)
		return
	}
	userInfo, err := c.auth.GetAccountInfoByUID(payment.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	if userInfo.Email == "" {
		log.Println("User don't have email")
		return
	}
	user, err := c.db.UserMustGet(payment.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	if user.Name == "" {
		user.Name = "Anonymous"
	}
	body := fmt.Sprintf(`สวัสดีครับคุณ %s,<br>
<br>
อีเมล์ฉบับนี้ยืนยันว่าท่านได้รับการอนุมัติการชำระเงินสำหรับหลักสูตร "%s" เสร็จสิ้น ท่านสามารถทำการ login เข้าสู่ Website Acourse แล้วเข้าเรียนหลักสูตร "%s" ได้ทันที<br>
<br>
รหัสการชำระเงิน: %s<br>
ชื่อหลักสูตร: %s<br>
จำนวนเงิน: %.2f บาท<br>
เวลาที่ทำการชำระเงิน: %s<br>
เวลาที่อนุมัติการชำระเงิน: %s<br>
ชื่อผู้ชำระเงิน: %s<br>
อีเมล์ผู้ชำระเงิน: %s<br>
`,
		user.Name,
		course.Title,
		course.Title,
		payment.ID,
		course.Title,
		payment.Price,
		payment.CreatedAt.Format(time.RFC822),
		payment.At.Format(time.RFC822),
		user.Name,
		userInfo.Email,
	)

	// 	if course.Type == model.CourseTypeVideo {
	// 		body += `----------------------
	// สนใจรับ Certificate หลังจบ Course <a href=''> เกี่ยวกับ ACertificate </a>

	// ท่านสามารถเพิ่มเงินจำนวน 600 บาท (จำนวน ชม. x 30 บาท  ) เพื่อได้รับ ACertificate หลังจากทำการส่งการบ้านครบถ้วน เพื่อใช้เป็นหลักฐานอ้างอิงกับ <a href=''>บริษัท Partner ของเรา</a>

	// ท่านสามารถสั่ง ACertificate ได้โดยโอนเงินจำนวน xxx บาท มาที่ บัญชีธนาคาร (ฝาก north เติมบัญชีตามหน้าเว็บ) แล้ว reply email นี้พร้อมแนบหลักฐานการโอนเงิน และเขียนว่า 'สั่งซื้อ certificate'
	// `
	// 	}

	body += `----------------------<br>
ขอบคุณที่ร่วมเรียนกับเราครับ<br>
Krissada Chalermsook (Oak)<br>
Founder/CEO Acourse.io<br>
https://acourse.io`

	err = SendMail(Email{
		To:      []string{userInfo.Email},
		Subject: fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", course.Title),
		Body:    body,
	})
	if err != nil {
		log.Println(err)
	}
}

// Reject runs reject action
func (c *PaymentController) Reject(ctx *app.PaymentRejectContext) error {
	role, err := c.db.RoleGet(ctx.CurrentUserID)
	if err != nil {
		return err
	}
	if !role.Admin {
		return e.ErrForbidden
	}

	payment, err := c.db.PaymentGet(ctx.PaymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return e.ErrNotFound
	}
	payment.Reject()

	err = c.db.PaymentSave(payment)
	if err != nil {
		return err
	}

	return nil
}
