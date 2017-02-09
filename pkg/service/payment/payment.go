package payment

import (
	"fmt"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/ds"
	"github.com/acoshift/go-firebase-admin"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new payment service
func New(client *ds.Client, user acourse.UserServiceClient, course acourse.CourseServiceClient, auth *admin.Auth, email acourse.EmailServiceClient) acourse.PaymentServiceServer {
	return &service{client, user, course, auth, email}
}

type service struct {
	client *ds.Client
	user   acourse.UserServiceClient
	course acourse.CourseServiceClient
	auth   *admin.Auth
	email  acourse.EmailServiceClient
}

func (s *service) validateUser(ctx context.Context) error {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	role, err := s.user.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
	if err != nil {
		return err
	}
	if !role.Admin {
		return grpc.Errorf(codes.PermissionDenied, "permission denied")
	}
	return nil
}

func (s *service) listPayments(ctx context.Context, opts ...ds.Query) (*acourse.PaymentsResponse, error) {
	err := s.validateUser(ctx)
	if err != nil {
		return nil, err
	}

	var payments []*paymentModel
	err = s.client.Query(ctx, kindPayment, &payments, opts...)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) || len(payments) == 0 {
		return &acourse.PaymentsResponse{}, nil
	}
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

	usersResp, err := s.user.GetUsers(ctx, &acourse.UserIDsRequest{UserIds: userIDs})
	if err != nil {
		return nil, err
	}
	users := usersResp.GetUsers()

	courses, err := s.course.GetCourses(ctx, &acourse.CourseIDsRequest{CourseIds: courseIDs})
	if err != nil {
		return nil, err
	}

	coursesTiny := make([]*acourse.CourseTiny, len(courses.Courses))
	for i, course := range courses.Courses {
		coursesTiny[i] = &acourse.CourseTiny{
			Id:    course.GetId(),
			Title: course.GetTitle(),
		}
	}

	return &acourse.PaymentsResponse{
		Payments: toPayments(payments),
		Users:    acourse.ToUsersTiny(users),
		Courses:  coursesTiny,
	}, nil
}

func (s *service) ListWaitingPayments(ctx context.Context, req *acourse.ListRequest) (*acourse.PaymentsResponse, error) {
	return s.listPayments(ctx, ds.Filter("Status =", string(statusWaiting)))
}

func (s *service) ListHistoryPayments(ctx context.Context, req *acourse.ListRequest) (*acourse.PaymentsResponse, error) {
	return s.listPayments(ctx)
}

func (s *service) ApprovePayments(ctx context.Context, req *acourse.PaymentIDsRequest) (*acourse.Empty, error) {
	err := s.validateUser(ctx)
	if err != nil {
		return nil, err
	}

	var payments []*paymentModel
	err = s.client.GetByStringIDs(ctx, kindPayment, req.GetPaymentIds(), &payments)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}

	enrolls := make([]*acourse.Enroll, 0, len(payments))
	for _, payment := range payments {
		payment.Approve()
		enrolls = append(enrolls, &acourse.Enroll{
			UserID:   payment.UserID,
			CourseID: payment.CourseID,
		})
	}

	_, err = s.course.CreateEnrolls(ctx, &acourse.EnrollsRequest{Enrolls: enrolls})
	if err != nil {
		return nil, err
	}

	err = s.client.SaveModels(ctx, "", payments)
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

	var payments []*paymentModel
	err = s.client.GetByStringIDs(ctx, kindPayment, req.GetPaymentIds(), &payments)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}

	for _, payment := range payments {
		payment.Reject()
	}

	err = s.client.SaveModels(ctx, "", payments)
	if err != nil {
		return nil, err
	}

	go s.sendRejectNotification(payments)

	return new(acourse.Empty), nil
}

func (s *service) processApprovedPayment(payment *paymentModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	courseResp, err := s.course.GetCourse(ctx, &acourse.CourseIDRequest{CourseId: payment.CourseID})
	if err != nil {
		return err
	}
	course := courseResp.GetCourse()
	userInfo, err := s.auth.GetUser(payment.UserID)
	if err != nil {
		return err
	}
	if userInfo.Email == "" {
		internal.WarningLogger.Println("User don't have email")
		return nil
	}
	user, err := s.user.GetUser(ctx, &acourse.UserIDRequest{UserId: payment.UserID})
	if err != nil {
		return err
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

	_, err = s.email.Send(ctx, &acourse.Email{
		To:      []string{userInfo.Email},
		Subject: fmt.Sprintf("ยืนยันการชำระเงิน หลักสูตร %s", course.Title),
		Body:    body,
	})
	return err
}

func (s *service) processRejectedPayment(payment *paymentModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	courseResp, err := s.course.GetCourse(ctx, &acourse.CourseIDRequest{CourseId: payment.CourseID})
	if err != nil {
		return err
	}
	course := courseResp.GetCourse()
	userInfo, err := s.auth.GetUser(payment.UserID)
	if err != nil {
		return err
	}
	if userInfo.Email == "" {
		internal.WarningLogger.Println("User don't have email")
		return nil
	}
	user, err := s.user.GetUser(ctx, &acourse.UserIDRequest{UserId: payment.UserID})
	if err != nil {
		return err
	}
	if user.Name == "" {
		user.Name = "Anonymous"
	}
	if len(course.Url) == 0 {
		course.Url = course.Id
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
		course.Url,
	)

	_, err = s.email.Send(ctx, &acourse.Email{
		To:      []string{userInfo.Email},
		Subject: fmt.Sprintf("คำขอเพื่อเรียนหลักสูตร %s ได้รับการปฏิเสธ", course.Title),
		Body:    body,
	})
	return err
}

func (s *service) sendApprovedNotification(payments []*paymentModel) {
	for _, payment := range payments {
		err := s.processApprovedPayment(payment)
		if err != nil {
			internal.ErrorLogger.Printf("PaymentService: sendApprovedNotification: %v", err)
		}
	}
}

func (s *service) sendRejectNotification(payments []*paymentModel) {
	for _, payment := range payments {
		err := s.processRejectedPayment(payment)
		if err != nil {
			internal.ErrorLogger.Printf("PaymentService: sendRejectNotification: %v", err)
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

	var x paymentModel
	err = s.client.GetByStringID(ctx, kindPayment, req.GetPaymentId(), &x)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}
	x.Price = req.GetPrice()

	err = s.client.SaveModel(ctx, "", &x)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}

func (s *service) FindPayment(ctx context.Context, req *acourse.PaymentFindRequest) (*acourse.Payment, error) {
	var x paymentModel
	err := s.client.QueryFirst(ctx, kindPayment, &x,
		ds.Filter("UserID =", req.GetUserId()),
		ds.Filter("CourseID =", req.GetCourseId()),
		ds.Filter("Status =", req.GetStatus()),
	)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}
	return toPayment(&x), nil
}

func (s *service) CreatePayment(ctx context.Context, req *acourse.Payment) (*acourse.Empty, error) {
	x := paymentModel{
		CourseID:      req.CourseId,
		UserID:        req.UserId,
		OriginalPrice: req.OriginalPrice,
		Price:         req.Price,
		Code:          req.Code,
		URL:           req.Url,
		Status:        status(req.Status),
	}

	err := s.client.SaveModel(ctx, kindPayment, &x)
	if err != nil {
		return nil, err
	}
	return new(acourse.Empty), nil
}

// StartNotification starts payment notification
func StartNotification(client *ds.Client, email acourse.EmailServiceClient) {
	ctx := context.Background()
	go func() {
		for {
			time.Sleep(6 * time.Hour)

			// check is any payments have status waiting
			internal.NoticeLogger.Println("PaymentService: Run Notification Payment")
			var payments []*paymentModel
			err := client.Query(ctx, kindPayment, &payments, ds.Filter("Status =", string(statusWaiting)))
			err = ds.IgnoreFieldMismatch(err)
			if err == nil && len(payments) > 0 {
				_, err = email.Send(ctx, &acourse.Email{
					To:      []string{"acoshift@gmail.com", "k.chalermsook@gmail.com"},
					Subject: "Acourse - Payment Received",
					Body:    fmt.Sprintf("%d payments pending", len(payments)),
				})
				if err != nil {
					internal.ErrorLogger.Printf("PaymentService: notification payment: %v", err)
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
