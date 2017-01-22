package payment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/app"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/go-firebase-admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Store is the payment store
type Store interface {
	PaymentList(context.Context, ...store.PaymentListOption) (model.Payments, error)
	PaymentGetMulti(context.Context, []string) (model.Payments, error)
	PaymentSaveMulti(context.Context, model.Payments) error
	UserGetMulti(context.Context, []string) (model.Users, error)
	UserMustGet(context.Context, string) (*model.User, error)
	CourseGet(context.Context, string) (*model.Course, error)
	CourseGetMulti(context.Context, []string) (model.Courses, error)
	EnrollSaveMulti(context.Context, []*model.Enroll) error
	RoleGet(context.Context, string) (*model.Role, error)
	PaymentGet(context.Context, string) (*model.Payment, error)
	PaymentSave(context.Context, *model.Payment) error
}

// New creates new payment service
func New(store Store, auth *admin.Auth, email acourse.EmailServiceClient) acourse.PaymentServiceServer {
	return &service{store, auth, email}
}

type service struct {
	store Store
	auth  *admin.Auth
	email acourse.EmailServiceClient
}

func (s *service) validateUser(ctx context.Context) error {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	role, err := s.store.RoleGet(ctx, userID)
	if err != nil {
		return err
	}
	if !role.Admin {
		return grpc.Errorf(codes.PermissionDenied, "permission denied")
	}
	return nil
}

func (s *service) listPayments(ctx context.Context, opts ...store.PaymentListOption) (*acourse.PaymentsResponse, error) {
	err := s.validateUser(ctx)
	if err != nil {
		return nil, err
	}

	payments, err := s.store.PaymentList(ctx, opts...)
	if err != nil {
		return nil, err
	}

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

	courses, err := s.store.CourseGetMulti(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	return &acourse.PaymentsResponse{
		Payments: acourse.ToPayments(payments),
		Users:    acourse.ToUsersTiny(users),
		Courses:  acourse.ToCoursesTiny(courses),
	}, nil
}

func (s *service) ListWaitingPayments(ctx context.Context, req *acourse.ListRequest) (*acourse.PaymentsResponse, error) {
	return s.listPayments(ctx, store.PaymentListOptionStatus(model.PaymentStatusWaiting))
}

func (s *service) ListHistoryPayments(ctx context.Context, req *acourse.ListRequest) (*acourse.PaymentsResponse, error) {
	return s.listPayments(ctx)
}

func (s *service) ApprovePayments(ctx context.Context, req *acourse.PaymentIDsRequest) (*acourse.Empty, error) {
	err := s.validateUser(ctx)
	if err != nil {
		return nil, err
	}

	payments, err := s.store.PaymentGetMulti(ctx, app.UniqueIDs(req.GetPaymentIds()))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	err = s.store.PaymentSaveMulti(ctx, payments)
	if err != nil {
		return nil, err
	}

	go s.sendApprovedNotification(payments)

	return new(acourse.Empty), nil
}

func (s *service) RejectPayments(ctx context.Context, req *acourse.PaymentIDsRequest) (*acourse.Empty, error) {
	err := s.validateUser(ctx)
	if err != nil {
		return nil, err
	}

	payments, err := s.store.PaymentGetMulti(ctx, app.UniqueIDs(req.GetPaymentIds()))
	if err != nil {
		return nil, err
	}

	for _, payment := range payments {
		payment.Reject()
	}

	err = s.store.PaymentSaveMulti(ctx, payments)
	if err != nil {
		return nil, err
	}

	go s.sendRejectNotification(payments)

	return new(acourse.Empty), nil
}

func (s *service) sendApprovedNotification(payments []*model.Payment) {
	ctx := context.Background()
	for _, payment := range payments {
		course, err := s.store.CourseGet(ctx, payment.CourseID)
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
		user, err := s.store.UserMustGet(ctx, payment.UserID)
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
			payment.ID(),
			course.Title,
			payment.Price,
			payment.CreatedAt.In(timeLocal).Format(time.RFC822),
			payment.At.In(timeLocal).Format(time.RFC822),
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

		_, err = s.email.Send(context.Background(), &acourse.Email{
			To:      []string{userInfo.Email},
			Subject: fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", course.Title),
			Body:    body,
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *service) sendRejectNotification(payments []*model.Payment) {
	ctx := context.Background()
	for _, payment := range payments {
		course, err := s.store.CourseGet(ctx, payment.CourseID)
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
		user, err := s.store.UserMustGet(ctx, payment.UserID)
		if err != nil {
			log.Println(err)
			continue
		}
		if user.Name == "" {
			user.Name = "Anonymous"
		}
		if course.URL == "" {
			course.URL = course.ID()
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
			user.Name,
			course.Title,
			payment.CreatedAt.In(timeLocal).Format(time.RFC822),
			course.URL,
		)

		_, err = s.email.Send(context.Background(), &acourse.Email{
			To:      []string{userInfo.Email},
			Subject: fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", course.Title),
			Body:    body,
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *service) UpdatePrice(ctx context.Context, req *acourse.PaymentUpdatePriceRequest) (*acourse.Empty, error) {
	err := s.validateUser(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetPaymentId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "payment id required")
	}

	payment, err := s.store.PaymentGet(ctx, req.GetPaymentId())
	if err != nil {
		return nil, err
	}
	payment.Price = req.GetPrice()

	err = s.store.PaymentSave(ctx, payment)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}

// StartNotification starts payment notification
func StartNotification(s Store, email acourse.EmailServiceClient) {
	go func() {
		for {
			time.Sleep(6 * time.Hour)

			// check is any payments have status waiting
			log.Println("Run Notification Payment")
			payments, err := s.PaymentList(context.Background(), store.PaymentListOptionStatus(model.PaymentStatusWaiting))
			if err == nil && len(payments) > 0 {
				_, err = email.Send(context.Background(), &acourse.Email{
					To:      []string{"acoshift@gmail.com", "k.chalermsook@gmail.com"},
					Subject: "Acourse - Payment Received",
					Body:    fmt.Sprintf("%d payments pending", len(payments)),
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

var timeLocal *time.Location

func init() {
	var err error
	timeLocal, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		timeLocal = time.UTC
	}
}
