package payment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/httperror"
)

// Store is the payment store
type Store interface {
	PaymentList(opts ...store.PaymentListOption) (model.Payments, error)
	PaymentGetMulti(context.Context, []string) (model.Payments, error)
	PaymentSaveMulti(context.Context, model.Payments) error
	UserGetMulti(context.Context, []string) (model.Users, error)
	UserMustGet(string) (*model.User, error)
	CourseGet(string) (*model.Course, error)
	CourseGetMulti(context.Context, []string) (model.Courses, error)
	EnrollSaveMulti(context.Context, []*model.Enroll) error
	RoleGet(string) (*model.Role, error)
}

// New creates new payment service
func New(store Store, auth *admin.Auth, email app.EmailService) app.PaymentService {
	return &service{store, auth, email}
}

type service struct {
	store Store
	auth  *admin.Auth
	email app.EmailService
}

func (s *service) ListPayments(ctx context.Context, req *app.PaymentListRequest) (*app.PaymentsReply, error) {
	currentUserID, ok := ctx.Value(app.KeyCurrentUserID).(string)
	if !ok {
		return nil, httperror.Unauthorized
	}
	role, err := s.store.RoleGet(currentUserID)
	if err != nil {
		return nil, err
	}
	if !role.Admin {
		return nil, httperror.Forbidden
	}

	var payments model.Payments
	if req.History != nil && *req.History {
		payments, err = s.store.PaymentList()
	} else {
		payments, err = s.store.PaymentList(store.PaymentListOptionStatus(model.PaymentStatusWaiting))
	}
	if err != nil {
		return nil, err
	}
	payments.SetView(model.PaymentViewDefault)

	userIDMap := map[string]bool{}
	courseIDMap := map[string]bool{}
	for _, payment := range payments {
		userIDMap[payment.UserID] = true
		courseIDMap[payment.CourseID] = true
	}
	userIDs := make([]string, 0, len(userIDMap))
	for userID := range userIDMap {
		userIDs = append(userIDs, userID)
	}
	courseIDs := make([]string, 0, len(courseIDMap))
	for courseID := range courseIDMap {
		courseIDs = append(courseIDs, courseID)
	}

	users, err := s.store.UserGetMulti(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	users.SetView(model.UserViewTiny)

	courses, err := s.store.CourseGetMulti(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	courses.SetView(model.CourseViewMini)

	return &app.PaymentsReply{
		Payments: payments,
		Users:    users,
		Courses:  courses,
	}, nil
}

func (s *service) ApprovePayments(ctx context.Context, req *app.IDsRequest) error {
	currentUserID, ok := ctx.Value(app.KeyCurrentUserID).(string)
	if !ok {
		return httperror.Unauthorized
	}
	role, err := s.store.RoleGet(currentUserID)
	if err != nil {
		return err
	}
	if !role.Admin {
		return httperror.Forbidden
	}

	payments, err := s.store.PaymentGetMulti(ctx, req.UniqueIDs())
	if err != nil {
		return err
	}

	enrolls := make([]*model.Enroll, 0, len(payments))
	for _, payment := range payments {
		payment.Approve()
		enrolls = append(enrolls, &model.Enroll{
			UserID:   payment.UserID,
			CourseID: payment.CourseID,
		})
	}

	err = s.store.EnrollSaveMulti(ctx, enrolls)
	if err != nil {
		return err
	}

	err = s.store.PaymentSaveMulti(ctx, payments)
	if err != nil {
		return err
	}

	go s.sendApprovedNotification(payments)

	return nil
}

func (s *service) RejectPayments(ctx context.Context, req *app.IDsRequest) error {
	currentUserID, ok := ctx.Value(app.KeyCurrentUserID).(string)
	if !ok {
		return httperror.Unauthorized
	}
	role, err := s.store.RoleGet(currentUserID)
	if err != nil {
		return err
	}
	if !role.Admin {
		return httperror.Forbidden
	}

	payments, err := s.store.PaymentGetMulti(ctx, req.UniqueIDs())
	if err != nil {
		return err
	}

	for _, payment := range payments {
		payment.Reject()
	}

	err = s.store.PaymentSaveMulti(ctx, payments)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) sendApprovedNotification(payments []*model.Payment) {
	for _, payment := range payments {
		course, err := s.store.CourseGet(payment.CourseID)
		if err != nil {
			log.Println(err)
			continue
		}
		userInfo, err := s.auth.GetAccountInfoByUID(payment.UserID)
		if err != nil {
			log.Println(err)
			continue
		}
		if userInfo.Email == "" {
			log.Println("User don't have email")
			continue
		}
		user, err := s.store.UserMustGet(payment.UserID)
		if err != nil {
			log.Println(err)
			continue
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

		err = s.email.SendEmail(nil, &app.EmailRequest{
			To:      []string{userInfo.Email},
			Subject: fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", course.Title),
			Body:    body,
		})
		if err != nil {
			log.Println(err)
		}
	}
}

// StartNotification starts payment notification
func StartNotification(s Store, email app.EmailService) {
	go func() {
		for {
			// check is payments have status waiting
			log.Println("Run Notification Payment")
			payments, err := s.PaymentList(store.PaymentListOptionStatus(model.PaymentStatusWaiting))
			if err == nil && len(payments) > 0 {
				err = email.SendEmail(nil, &app.EmailRequest{
					To:      []string{"acoshift@gmail.com", "k.chalermsook@gmail.com"},
					Subject: "Acourse - Payment Received",
					Body:    fmt.Sprintf("%d payments pending", len(payments)),
				})
				if err != nil {
					log.Println(err)
				}
			}
			time.Sleep(3 * time.Hour)
		}
	}()
}
