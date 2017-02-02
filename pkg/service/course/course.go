package course

import (
	"context"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/ds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new course service
func New(client *ds.Client, user acourse.UserServiceClient, payment acourse.PaymentServiceClient) acourse.CourseServiceServer {
	return &service{client, user, payment}
}

type service struct {
	client  *ds.Client
	user    acourse.UserServiceClient
	payment acourse.PaymentServiceClient
}

func (s *service) listCourses(ctx context.Context, qs ...ds.Query) (*acourse.CoursesResponse, error) {
	var courses courseModels
	err := s.client.Query(ctx, kindCourse, &courses, qs...)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	// get owners
	userIDMap := map[string]bool{}
	for _, x := range courses {
		userIDMap[x.Owner] = true
	}
	userIDs := make([]string, 0, len(userIDMap))
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}
	usersResp, err := s.user.GetUsers(ctx, &acourse.UserIDsRequest{UserIds: userIDs})
	if err != nil {
		return nil, err
	}
	users := usersResp.GetUsers()

	enrollCounts := make([]*acourse.EnrollCount, 0, len(courses))
	enrollResult := make(chan *acourse.EnrollCount)
	for i, x := range courses {
		go func(i int, x *courseModel) {
			c, err := s.countEnroll(ctx, x.ID())
			if err != nil {
				panic(err)
			}
			enrollResult <- &acourse.EnrollCount{
				CourseId: x.ID(),
				Count:    int32(c),
			}
		}(i, x)
	}
	for range courses {
		enrollCounts = append(enrollCounts, <-enrollResult)
	}
	return &acourse.CoursesResponse{
		Courses:      acourse.ToCoursesSmall(toCourses(courses)),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}

func (s *service) ListPublicCourses(ctx context.Context, req *acourse.ListRequest) (*acourse.CoursesResponse, error) {
	return s.listCourses(ctx, ds.Filter("Options.Public =", true))
}

func (s *service) ListCourses(ctx context.Context, req *acourse.ListRequest) (*acourse.CoursesResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	role, err := s.user.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if !role.Admin {
		return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
	}
	return s.listCourses(ctx)
}

func (s *service) ListOwnCourses(ctx context.Context, req *acourse.UserIDRequest) (*acourse.CoursesResponse, error) {
	userID := internal.GetUserID(ctx)

	if len(req.UserId) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid user id")
	}

	qs := make([]ds.Query, 0, 3)
	qs = append(qs, ds.Filter("Owner =", req.UserId))

	// if not sign in, get only public courses
	if len(userID) == 0 {
		qs = append(qs, ds.Filter("Options.Public =", true))
	}

	return s.listCourses(ctx, qs...)
}

func (s *service) ListEnrolledCourses(ctx context.Context, req *acourse.UserIDRequest) (*acourse.CoursesResponse, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	if req.GetUserId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid user id")
	}

	// only admin allow for get other user enrolled courses
	if req.GetUserId() != userID {
		role, err := s.user.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
		if err != nil {
			return nil, err
		}
		if !role.Admin {
			return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
		}
	}

	enrolls, err := s.listEnrollByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(enrolls))
	for i, e := range enrolls {
		ids[i] = e.CourseID
	}
	var courses courseModels
	err = s.client.GetByStringIDs(ctx, kindCourse, ids, &courses)
	err = ds.IgnoreFieldMismatch(err)

	// get owners
	userIDMap := map[string]bool{}
	for _, course := range courses {
		userIDMap[course.Owner] = true
	}
	userIDs := make([]string, 0, len(userIDMap))
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}
	usersResp, err := s.user.GetUsers(ctx, &acourse.UserIDsRequest{UserIds: userIDs})
	if err != nil {
		return nil, err
	}
	users := usersResp.GetUsers()

	enrollCounts := make([]*acourse.EnrollCount, len(courses))
	for i, course := range courses {
		c, err := s.countEnroll(ctx, course.ID())
		if err != nil {
			return nil, err
		}
		enrollCounts[i] = &acourse.EnrollCount{
			CourseId: course.ID(),
			Count:    int32(c),
		}
	}
	return &acourse.CoursesResponse{
		Courses:      acourse.ToCoursesSmall(toCourses(courses)),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}

func (s *service) GetCourse(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.CourseResponse, error) {
	userID := internal.GetUserID(ctx)

	var course courseModel

	// try get by id first
	err := s.client.GetByStringID(ctx, kindCourse, req.GetCourseId(), &course)
	if ds.NotFound(err) {
		err = s.client.QueryFirst(ctx, kindCourse, &course, ds.Filter("URL =", req.GetCourseId()))
	}
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, errCourseNotFound
	}
	if err != nil {
		return nil, err
	}

	// get course owner
	owner, err := s.user.GetUser(ctx, &acourse.UserIDRequest{UserId: course.Owner})
	if err != nil {
		return nil, err
	}

	// check is user enrolled
	enroll, err := s.FindEnroll(ctx, &acourse.EnrollFindRequest{UserId: userID, CourseId: course.ID()})
	if grpc.Code(err) == codes.NotFound {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if enroll != nil || course.Owner == userID {
		var attend *attendModel
		attend, err = s.findAttend(ctx, userID, course.ID())
		if grpc.Code(err) == codes.NotFound {
			err = nil
		}
		if err != nil {
			return nil, err
		}

		return &acourse.CourseResponse{
			Course:   toCourse(&course),
			User:     acourse.ToUserTiny(owner),
			Enrolled: enroll != nil,
			Owned:    course.Owner == userID,
			Attended: attend != nil,
		}, nil
	}

	// check waiting payment
	var payment *acourse.Payment
	if userID != "" {
		payment, err = s.payment.FindPayment(ctx, &acourse.PaymentFindRequest{
			UserId:   userID,
			CourseId: course.ID(),
			Status:   "waiting",
		})
		if err != nil && grpc.Code(err) != codes.NotFound {
			return nil, err
		}
	}

	role, err := s.user.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	if role.Admin {
		return &acourse.CourseResponse{
			Course:   toCourse(&course),
			User:     acourse.ToUserTiny(owner),
			Enrolled: enroll != nil,
			Purchase: payment != nil,
		}, nil
	}

	// filter out private fields
	course = courseModel{
		StringIDModel:    course.StringIDModel,
		StampModel:       course.StampModel,
		Title:            course.Title,
		ShortDescription: course.ShortDescription,
		Description:      course.Description,
		Photo:            course.Photo,
		Owner:            course.Owner,
		Start:            course.Start,
		URL:              course.URL,
		Type:             course.Type,
		Price:            course.Price,
		DiscountedPrice:  course.DiscountedPrice,
		Options: courseOption{
			Public:   course.Options.Public,
			Discount: course.Options.Discount,
			Enroll:   course.Options.Enroll,
		},
		EnrollDetail: course.EnrollDetail,
	}

	return &acourse.CourseResponse{
		Course:   toCourse(&course),
		User:     acourse.ToUserTiny(owner),
		Purchase: payment != nil,
	}, nil
}

func (s *service) GetCourses(ctx context.Context, req *acourse.CourseIDsRequest) (*acourse.CoursesResponse, error) {
	var courses []*courseModel
	err := s.client.GetByStringIDs(ctx, kindCourse, req.GetCourseIds(), &courses)
	err = ds.IgnoreFieldMismatch(err)
	if err != nil {
		return nil, err
	}
	return &acourse.CoursesResponse{Courses: acourse.ToCoursesSmall(toCourses(courses))}, nil
}

func (s *service) CreateCourse(ctx context.Context, req *acourse.Course) (*acourse.Course, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	role, err := s.user.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if !role.Instructor {
		return nil, grpc.Errorf(codes.PermissionDenied, "don't have permission to create course")
	}

	course := &courseModel{
		Title:            req.GetTitle(),
		ShortDescription: req.GetShortDescription(),
		Description:      req.GetDescription(),
		Photo:            req.GetPhoto(),
		Video:            req.GetVideo(),
		Owner:            userID,
		Options: courseOption{
			Assignment: req.GetOptions().GetAssignment(),
		},
	}
	course.Start, _ = time.Parse(time.RFC3339, req.GetStart())
	course.Contents = make(courseContents, len(req.GetContents()))
	for i, c := range req.GetContents() {
		course.Contents[i] = courseContent{
			Title:       c.GetTitle(),
			Description: c.GetDescription(),
			Video:       c.GetVideo(),
			DownloadURL: c.GetDownloadURL(),
		}
	}

	err = s.client.SaveModel(ctx, kindCourse, course)
	if err != nil {
		return nil, err
	}

	return toCourse(course), nil
}

func (s *service) UpdateCourse(ctx context.Context, req *acourse.Course) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var course courseModel
	err := s.client.GetByStringID(ctx, kindCourse, req.GetId(), &course)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, errCourseNotFound
	}
	if err != nil {
		return nil, err
	}
	role, err := s.user.GetRole(ctx, &acourse.UserIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if course.Owner != userID && !role.Admin {
		return nil, grpc.Errorf(codes.PermissionDenied, "don't have permission to update this course")
	}

	// merge course with request
	course.Title = req.GetTitle()
	course.ShortDescription = req.GetShortDescription()
	course.Description = req.GetDescription()
	course.Photo = req.GetPhoto()
	course.Start, _ = time.Parse(time.RFC3339, req.GetStart())
	course.Video = req.GetVideo()
	course.Contents = make(courseContents, len(req.GetContents()))
	for i, c := range req.GetContents() {
		course.Contents[i] = courseContent{
			Title:       c.GetTitle(),
			Description: c.GetDescription(),
			Video:       c.GetVideo(),
			DownloadURL: c.GetDownloadURL(),
		}
	}
	course.Options.Assignment = req.GetOptions().GetAssignment()

	err = s.client.SaveModel(ctx, "", &course)
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}

func (s *service) EnrollCourse(ctx context.Context, req *acourse.EnrollRequest) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	if req.CourseId == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "course id required")
	}

	var course courseModel
	err := s.client.GetByStringID(ctx, kindCourse, req.GetCourseId(), &course)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, errCourseNotFound
	}
	if err != nil {
		return nil, err
	}

	// owner can not enroll
	if course.Owner == userID {
		return nil, grpc.Errorf(codes.PermissionDenied, "owner can not enroll their own course")
	}

	// check is user already enroll
	enroll, err := s.FindEnroll(ctx, &acourse.EnrollFindRequest{UserId: userID, CourseId: req.GetCourseId()})
	if grpc.Code(err) == codes.NotFound {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if enroll != nil {
		// user already enroll
		return nil, grpc.Errorf(codes.AlreadyExists, "already enroll")
	}

	// check is user already send waiting payment
	_, err = s.payment.FindPayment(ctx, &acourse.PaymentFindRequest{
		UserId:   userID,
		CourseId: req.CourseId,
		Status:   "waiting",
	})
	if err != nil && grpc.Code(err) != codes.NotFound {
		return nil, err
	}
	if err == nil {
		// user already send payment
		return nil, grpc.Errorf(codes.FailedPrecondition, "wait admin to review your current payment before send another payment for this course")
	}

	// calculate price
	originalPrice := course.Price
	if course.Options.Discount {
		originalPrice = course.DiscountedPrice
	}
	// TODO: calculate code

	// auto enroll if course free
	if originalPrice == 0.0 {
		enroll := &enrollModel{
			UserID:   userID,
			CourseID: req.CourseId,
		}
		err = s.saveEnroll(ctx, enroll)
		if err != nil {
			return nil, err
		}
		return new(acourse.Empty), nil
	}

	// create payment
	_, err = s.payment.CreatePayment(ctx, &acourse.Payment{
		CourseId:      req.CourseId,
		UserId:        userID,
		OriginalPrice: originalPrice,
		Price:         req.Price,
		Code:          req.Code,
		Url:           req.Url,
		Status:        "waiting",
	})
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}

func (s *service) AttendCourse(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var course courseModel
	err := s.client.GetByStringID(ctx, kindCourse, req.GetCourseId(), &course)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, errCourseNotFound
	}
	if err != nil {
		return nil, err
	}

	// user must enrolled in this course
	enroll, err := s.FindEnroll(ctx, &acourse.EnrollFindRequest{UserId: userID, CourseId: course.ID()})
	if grpc.Code(err) == codes.NotFound {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if enroll == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user must enroll first")
	}

	// check is user already attend
	attend, err := s.findAttend(ctx, userID, course.ID())
	if grpc.Code(err) == codes.NotFound {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if attend != nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "already attend in last 8 hr")
	}

	err = s.saveAttend(ctx, &attendModel{UserID: userID, CourseID: course.ID()})
	if err != nil {
		return nil, err
	}

	return new(acourse.Empty), nil
}

func (s *service) changeAttend(ctx context.Context, req *acourse.CourseIDRequest, value bool) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	var course courseModel
	err := s.client.GetByStringID(ctx, kindCourse, req.GetCourseId(), &course)
	err = ds.IgnoreFieldMismatch(err)
	if ds.NotFound(err) {
		return nil, errCourseNotFound
	}
	if err != nil {
		return nil, err
	}

	if course.Owner != userID {
		return nil, grpc.Errorf(codes.PermissionDenied, "don't have permission to change attend for this course")
	}

	course.Options.Attend = value

	err = s.client.SaveModel(ctx, "", &course)
	if err != nil {
		return nil, err
	}

	// TODO: notify users

	return new(acourse.Empty), nil
}

func (s *service) OpenAttend(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.Empty, error) {
	return s.changeAttend(ctx, req, true)
}

func (s *service) CloseAttend(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.Empty, error) {
	return s.changeAttend(ctx, req, false)
}

func (s *service) saveCourse(ctx context.Context, course *courseModel) error {
	var err error
	// Check duplicate URL
	if len(course.URL) > 0 {
		var t courseModel
		err = s.client.QueryFirst(ctx, kindCourse, &t, ds.Filter("URL =", course.URL))
		err = ds.IgnoreFieldMismatch(err)
		if !ds.NotFound(err) && t.ID() != course.ID() {
			return errCourseURLExists
		}
	}

	err = s.client.SaveModel(ctx, kindCourse, course)
	if err != nil {
		return err
	}

	return nil
}
