package course

import (
	"context"
	"time"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/acourse/pkg/model"
	"github.com/acoshift/acourse/pkg/store"
	_context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// New creates new course service
func New(store Store) acourse.CourseServiceServer {
	return &service{store}
}

// Store is the store interface for course service
type Store interface {
	CourseList(opts ...store.CourseListOption) (model.Courses, error)
	UserGetMulti(context.Context, []string) (model.Users, error)
	EnrollCourseCount(string) (int, error)
	RoleGet(string) (*model.Role, error)
	EnrollListByUserID(string) (model.Enrolls, error)
	CourseGetAllByIDs([]string) (model.Courses, error)
	CourseGet(string) (*model.Course, error)
	EnrollFind(string, string) (*model.Enroll, error)
	PaymentFind(string, string, model.PaymentStatus) (*model.Payment, error)
	EnrollSave(*model.Enroll) error
	PaymentSave(*model.Payment) error
	CourseSave(*model.Course) error
	CourseFind(string) (*model.Course, error)
	UserMustGet(string) (*model.User, error)
}

type service struct {
	store Store
}

func (s *service) listCourses(ctx _context.Context, opts ...store.CourseListOption) (*acourse.CoursesResponse, error) {
	courses, err := s.store.CourseList(opts...)
	if err != nil {
		return nil, err
	}
	// get owners
	userIDMap := map[string]bool{}
	for _, course := range courses {
		userIDMap[course.Owner] = true
	}
	userIDs := make([]string, 0, len(userIDMap))
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.store.UserGetMulti(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	enrollCounts := make([]*acourse.EnrollCount, len(courses))
	for i, course := range courses {
		c, err := s.store.EnrollCourseCount(course.ID)
		if err != nil {
			return nil, err
		}
		enrollCounts[i] = &acourse.EnrollCount{
			CourseId: course.ID,
			Count:    int32(c),
		}
	}
	return &acourse.CoursesResponse{
		Courses:      acourse.ToCoursesSmall(courses),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}

func (s *service) ListPublicCourses(ctx _context.Context, req *acourse.ListRequest) (*acourse.CoursesResponse, error) {
	return s.listCourses(ctx, store.CourseListOptionPublic(true))
}

func (s *service) ListCourses(ctx _context.Context, req *acourse.ListRequest) (*acourse.CoursesResponse, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	role, err := s.store.RoleGet(userID)
	if err != nil {
		return nil, err
	}
	if !role.Admin {
		return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
	}
	return s.listCourses(ctx)
}

func (s *service) ListOwnCourses(ctx _context.Context, req *acourse.UserIDRequest) (*acourse.CoursesResponse, error) {
	userID, _ := ctx.Value(acourse.KeyUserID).(string)

	if req.GetUserId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid user id")
	}

	opt := make([]store.CourseListOption, 0, 3)
	opt = append(opt, store.CourseListOptionOwner(req.GetUserId()))

	// if not sign in, get only public courses
	if userID == "" {
		opt = append(opt, store.CourseListOptionPublic(true))
	}

	return s.listCourses(ctx, opt...)
}

func (s *service) ListEnrolledCourses(ctx _context.Context, req *acourse.UserIDRequest) (*acourse.CoursesResponse, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	if req.GetUserId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid user id")
	}

	// only admin allow for get other user enrolled courses
	if req.GetUserId() != userID {
		role, err := s.store.RoleGet(userID)
		if err != nil {
			return nil, err
		}
		if !role.Admin {
			return nil, grpc.Errorf(codes.PermissionDenied, "permission denied")
		}
	}

	enrolls, err := s.store.EnrollListByUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(enrolls))
	for i, e := range enrolls {
		ids[i] = e.CourseID
	}
	courses, err := s.store.CourseGetAllByIDs(ids)

	// get owners
	userIDMap := map[string]bool{}
	for _, course := range courses {
		userIDMap[course.Owner] = true
	}
	userIDs := make([]string, 0, len(userIDMap))
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.store.UserGetMulti(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	enrollCounts := make([]*acourse.EnrollCount, len(courses))
	for i, course := range courses {
		c, err := s.store.EnrollCourseCount(course.ID)
		if err != nil {
			return nil, err
		}
		enrollCounts[i] = &acourse.EnrollCount{
			CourseId: course.ID,
			Count:    int32(c),
		}
	}
	return &acourse.CoursesResponse{
		Courses:      acourse.ToCoursesSmall(courses),
		Users:        acourse.ToUsersTiny(users),
		EnrollCounts: enrollCounts,
	}, nil
}

func (s *service) GetCourse(ctx _context.Context, req *acourse.CourseIDRequest) (*acourse.CourseResponse, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	// try get by id first
	course, err := s.store.CourseGet(req.GetCourseId())
	if err != nil {
		return nil, err
	}
	// try get by url
	if course == nil {
		course, err = s.store.CourseFind(req.GetCourseId())
		if err != nil {
			return nil, err
		}
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
	}

	// get course owner
	owner, err := s.store.UserMustGet(course.Owner)
	if err != nil {
		return nil, err
	}

	role, err := s.store.RoleGet(userID)
	if err != nil {
		return nil, err
	}

	// check is user enrolled
	enroll, err := s.store.EnrollFind(userID, course.ID)
	if err != nil {
		return nil, err
	}
	if enroll != nil || course.Owner == userID || role.Admin {
		return &acourse.CourseResponse{
			Course:   acourse.ToCourse(course),
			User:     acourse.ToUserTiny(owner),
			Enrolled: enroll != nil,
			Owned:    course.Owner == userID,
		}, nil
	}

	// filter out private fields
	course = &model.Course{
		Base:             course.Base,
		Stampable:        course.Stampable,
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
		},
		EnrollDetail: course.EnrollDetail,
	}

	// check waiting payment
	payment, err := s.store.PaymentFind(userID, course.ID, model.PaymentStatusWaiting)
	if err != nil {
		return nil, err
	}

	return &acourse.CourseResponse{
		Course:   acourse.ToCourse(course),
		User:     acourse.ToUserTiny(owner),
		Purchase: payment != nil,
	}, nil
}

func (s *service) CreateCourse(ctx _context.Context, req *acourse.Course) (*acourse.Course, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}
	role, err := s.store.RoleGet(userID)
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
			Attend:     req.GetOptions().GetAttend(),
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

	err = s.store.CourseSave(course)
	if err != nil {
		return nil, err
	}

	return acourse.ToCourse(course), nil
}

func (s *service) UpdateCourse(ctx _context.Context, req *acourse.Course) (*acourse.Empty, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	course, err := s.store.CourseGet(req.GetId())
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, grpc.Errorf(codes.NotFound, "course not found")
	}
	role, err := s.store.RoleGet(userID)
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
	course.Options.Attend = req.GetOptions().GetAttend()
	course.Options.Assignment = req.GetOptions().GetAssignment()

	err = s.store.CourseSave(course)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) EnrollCourse(ctx _context.Context, req *acourse.EnrollRequest) (*acourse.Empty, error) {
	userID, ok := ctx.Value(acourse.KeyUserID).(string)
	if !ok || userID == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "authorization required")
	}

	if req.GetCourseId() == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "course id required")
	}

	course, err := s.store.CourseGet(req.GetCourseId())
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
	enroll, err := s.store.EnrollFind(userID, req.GetCourseId())
	if err != nil {
		return nil, err
	}
	if enroll != nil {
		// user already enroll
		return nil, grpc.Errorf(codes.AlreadyExists, "already enroll")
	}

	// check is user already send waiting payment
	payment, err := s.store.PaymentFind(userID, req.GetCourseId(), model.PaymentStatusWaiting)
	if err != nil {
		return nil, err
	}
	if payment != nil {
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
		enroll = &model.Enroll{
			UserID:   userID,
			CourseID: req.GetCourseId(),
		}
		err = s.store.EnrollSave(enroll)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// create payment
	payment = &model.Payment{
		CourseID:      req.GetCourseId(),
		UserID:        userID,
		OriginalPrice: originalPrice,
		Price:         req.GetPrice(),
		Code:          req.GetCode(),
		URL:           req.GetUrl(),
		Status:        model.PaymentStatusWaiting,
	}

	err = s.store.PaymentSave(payment)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
