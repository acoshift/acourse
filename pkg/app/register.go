package app

import (
	"net/http"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/mount"
	"google.golang.org/grpc/metadata"
)

var m = mount.New(mount.Config{
	Binder:         bindJSON,
	SuccessHandler: handleOK,
	ErrorHandler:   handleError,
})

func serviceContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md := metadata.MD{}
		header := r.Header
		if v := header.Get("Authorization"); v != "" {
			md = metadata.Join(md, metadata.Pairs("authorization", v))
		}
		h.ServeHTTP(w, r.WithContext(metadata.NewContext(r.Context(), md)))
	})
}

// RegisterUserServiceClient registers a User service client to http server
func RegisterUserServiceClient(mux *http.ServeMux, s acourse.UserServiceClient) {
	sv := "/acourse.UserService"

	mux.Handle(sv+"/GetUser", serviceContext(m.Handler(s.GetUser)))
	mux.Handle(sv+"/GetUsers", serviceContext(m.Handler(s.GetUsers)))
	mux.Handle(sv+"/GetMe", serviceContext(m.Handler(s.GetMe)))
	mux.Handle(sv+"/UpdateMe", serviceContext(m.Handler(s.UpdateMe)))
}

// RegisterPaymentServiceClient registers a Payment service client to http server
func RegisterPaymentServiceClient(mux *http.ServeMux, s acourse.PaymentServiceClient) {
	sv := "/acourse.PaymentService"

	mux.Handle(sv+"/ListWaitingPayments", serviceContext(m.Handler(s.ListWaitingPayments)))
	mux.Handle(sv+"/ListHistoryPayments", serviceContext(m.Handler(s.ListHistoryPayments)))
	mux.Handle(sv+"/ApprovePayments", serviceContext(m.Handler(s.ApprovePayments)))
	mux.Handle(sv+"/RejectPayments", serviceContext(m.Handler(s.RejectPayments)))
	mux.Handle(sv+"/UpdatePrice", serviceContext(m.Handler(s.UpdatePrice)))
}

// RegisterCourseServiceClient registers a Course service client to http server
func RegisterCourseServiceClient(mux *http.ServeMux, s acourse.CourseServiceClient) {
	sv := "/acourse.CourseService"

	mux.Handle(sv+"/ListCourses", serviceContext(m.Handler(s.ListCourses)))
	mux.Handle(sv+"/ListPublicCourses", serviceContext(m.Handler(s.ListPublicCourses)))
	mux.Handle(sv+"/ListOwnCourses", serviceContext(m.Handler(s.ListOwnCourses)))
	mux.Handle(sv+"/ListEnrolledCourses", serviceContext(m.Handler(s.ListEnrolledCourses)))
	mux.Handle(sv+"/GetCourse", serviceContext(m.Handler(s.GetCourse)))
	mux.Handle(sv+"/CreateCourse", serviceContext(m.Handler(s.CreateCourse)))
	mux.Handle(sv+"/UpdateCourse", serviceContext(m.Handler(s.UpdateCourse)))
	mux.Handle(sv+"/EnrollCourse", serviceContext(m.Handler(s.EnrollCourse)))
	mux.Handle(sv+"/AttendCourse", serviceContext(m.Handler(s.AttendCourse)))
	mux.Handle(sv+"/OpenAttend", serviceContext(m.Handler(s.OpenAttend)))
	mux.Handle(sv+"/CloseAttend", serviceContext(m.Handler(s.CloseAttend)))
}

// RegisterAssignmentServiceClient registers s Assignment service client to http server
func RegisterAssignmentServiceClient(mux *http.ServeMux, s acourse.AssignmentServiceClient) {
	sv := "/acourse.AssignmentService"

	mux.Handle(sv+"/CreateAssignment", serviceContext(m.Handler(s.CreateAssignment)))
	mux.Handle(sv+"/UpdateAssignment", serviceContext(m.Handler(s.UpdateAssignment)))
	mux.Handle(sv+"/OpenAssignment", serviceContext(m.Handler(s.OpenAssignment)))
	mux.Handle(sv+"/CloseAssignment", serviceContext(m.Handler(s.CloseAssignment)))
	mux.Handle(sv+"/ListAssignments", serviceContext(m.Handler(s.ListAssignments)))
	mux.Handle(sv+"/DeleteAssignment", serviceContext(m.Handler(s.DeleteAssignment)))
	mux.Handle(sv+"/SubmitUserAssignment", serviceContext(m.Handler(s.SubmitUserAssignment)))
	mux.Handle(sv+"/DeleteUserAssignment", serviceContext(m.Handler(s.DeleteUserAssignment)))
	mux.Handle(sv+"/GetUserAssignments", serviceContext(m.Handler(s.GetUserAssignments)))
	mux.Handle(sv+"/ListUserAssignments", serviceContext(m.Handler(s.ListUserAssignments)))
}
