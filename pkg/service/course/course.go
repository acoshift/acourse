package course

import (
	"context"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/internal"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	"github.com/acoshift/ds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new course service
func New(store Store, client *ds.Client, user acourse.UserServiceClient, payment acourse.PaymentServiceClient) acourse.CourseServiceServer {
	return &service{store, client, user, payment}
}

// Store is the store interface for course service
type Store interface {
	CourseList(context.Context, ...store.CourseListOption) (model.Courses, error)
	CourseGetAllByIDs(context.Context, []string) (model.Courses, error)
	CourseGet(context.Context, string) (*model.Course, error)
	CourseSave(context.Context, *model.Course) error
	CourseFind(context.Context, string) (*model.Course, error)
}

type service struct {
	store   Store
	client  *ds.Client
	user    acourse.UserServiceClient
	payment acourse.PaymentServiceClient
}

func (s *service) listCourses(ctx context.Context, opts ...store.CourseListOption) (*acourse.CoursesResponse, error) {
	courses, err := s.store.CourseList(ctx, opts...)
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
		go func(i int, x *model.Course) {
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
		Courses:      acourse.ToCoursesSmall(courses),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}

func (s *service) ListPublicCourses(ctx context.Context, req *acourse.ListRequest) (*acourse.CoursesResponse, error) {
	return s.listCourses(ctx, store.CourseListOptionPublic(true))
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

	opt := make([]store.CourseListOption, 0, 3)
	opt = append(opt, store.CourseListOptionOwner(req.UserId))

	// if not sign in, get only public courses
	if len(userID) == 0 {
		opt = append(opt, store.CourseListOptionPublic(true))
	}

	return s.listCourses(ctx, opt...)
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
	courses, err := s.store.CourseGetAllByIDs(ctx, ids)

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
		Courses:      acourse.ToCoursesSmall(courses),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}

func (s *service) GetCourse(ctx context.Context, req *acourse.CourseIDRequest) (*acourse.CourseResponse, error) {
	userID := internal.GetUserID(ctx)

	// try get by id first
	course, err := s.store.CourseGet(ctx, req.CourseId)
	if err != nil {
		return nil, err
	}
	// try get by url
	if course == nil {
		course, err = s.store.CourseFind(ctx, req.CourseId)
		if err != nil {
			return nil, err
		}
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
	}

	// get course owner
	owner, err := s.user.GetUser(ctx, &acourse.UserIDRequest{UserId: course.Owner})
	if err != nil {
		return nil, err
	}

	// check is user enrolled
	enroll, err := s.FindEnroll(ctx, &acourse.EnrollFindRequest{UserId: userID, CourseId: course.ID()})
	if err != nil {
		return nil, err
	}
	if enroll != nil || course.Owner == userID {
		var attend *attendModel
		attend, err = s.findAttend(ctx, userID, course.ID())
		if err != nil {
			return nil, err
		}

		return &acourse.CourseResponse{
			Course:   acourse.ToCourse(course),
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
			Course:   acourse.ToCourse(course),
			User:     acourse.ToUserTiny(owner),
			Enrolled: enroll != nil,
			Purchase: payment != nil,
		}, nil
	}

	// filter out private fields
	course = &model.Course{
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
		Options: model.CourseOption{
			Public:   course.Options.Public,
			Discount: course.Options.Discount,
			Enroll:   course.Options.Enroll,
		},
		EnrollDetail: course.EnrollDetail,
	}

	return &acourse.CourseResponse{
		Course:   acourse.ToCourse(course),
		User:     acourse.ToUserTiny(owner),
		Purchase: payment != nil,
	}, nil
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

	course := &model.Course{
		Title:            req.GetTitle(),
		ShortDescription: req.GetShortDescription(),
		Description:      req.GetDescription(),
		Photo:            req.GetPhoto(),
		Video:            req.GetVideo(),
		Owner:            userID,
		Options: model.CourseOption{
			Assignment: req.GetOptions().GetAssignment(),
		},
	}
	course.Start, _ = time.Parse(time.RFC3339, req.GetStart())
	course.Contents = make(model.CourseContents, len(req.GetContents()))
	for i, c := range req.GetContents() {
		course.Contents[i] = model.CourseContent{
			Title:       c.GetTitle(),
			Description: c.GetDescription(),
			Video:       c.GetVideo(),
			DownloadURL: c.GetDownloadURL(),
		}
	}

	err = s.store.CourseSave(ctx, course)
	if err != nil {
		return nil, err
	}

	return acourse.ToCourse(course), nil
}

func (s *service) UpdateCourse(ctx context.Context, req *acourse.Course) (*acourse.Empty, error) {
	userID := internal.GetUserID(ctx)
	if userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	course, err := s.store.CourseGet(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
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
	course.Contents = make(model.CourseContents, len(req.GetContents()))
	for i, c := range req.GetContents() {
		course.Contents[i] = model.CourseContent{
			Title:       c.GetTitle(),
			Description: c.GetDescription(),
			Video:       c.GetVideo(),
			DownloadURL: c.GetDownloadURL(),
		}
	}
	course.Options.Assignment = req.GetOptions().GetAssignment()

	err = s.store.CourseSave(ctx, course)
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

	course, err := s.store.CourseGet(ctx, req.CourseId)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
	}

	// owner can not enroll
	if course.Owner == userID {
		return nil, grpc.Errorf(codes.PermissionDenied, "owner can not enroll their own course")
	}

	// check is user already enroll
	enroll, err := s.FindEnroll(ctx, &acourse.EnrollFindRequest{UserId: userID, CourseId: req.GetCourseId()})
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

	course, err := s.store.CourseGet(ctx, req.GetCourseId())
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
	}

	// user must enrolled in this course
	enroll, err := s.FindEnroll(ctx, &acourse.EnrollFindRequest{UserId: userID, CourseId: course.ID()})
	if err != nil {
		return nil, err
	}
	if enroll == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "user must enroll first")
	}

	// check is user already attend
	attend, err := s.findAttend(ctx, userID, course.ID())
	if err != nil {
		return nil, err
	}
	if attend != nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "already attend in last 6 hr")
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

	course, err := s.store.CourseGet(ctx, req.GetCourseId())
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
	}

	if course.Owner != userID {
		return nil, grpc.Errorf(codes.PermissionDenied, "don't have permission to change attend for this course")
	}

	course.Options.Attend = value

	err = s.store.CourseSave(ctx, course)
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
