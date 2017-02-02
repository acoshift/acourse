package app

import (
	"context"
	"net/http"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/httperror"
	"google.golang.org/grpc/metadata"
)

func makeServiceContext(r *http.Request) context.Context {
	md := metadata.MD{}
	header := r.Header
	if v := header.Get("Authorization"); v != "" {
		md = metadata.Join(md, metadata.Pairs("authorization", v))
	}
	return metadata.NewContext(r.Context(), md)
}

// RegisterUserServiceClient registers a User service client to http server
func RegisterUserServiceClient(mux *http.ServeMux, s acourse.UserServiceClient) {
	sv := "/acourse.UserService"

	mux.HandleFunc(sv+"/GetUser", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.UserIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetUser(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/GetUsers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.UserIDsRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetUsers(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/GetMe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		res, err := s.GetMe(makeServiceContext(r), new(acourse.Empty))
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/UpdateMe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.User)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.UpdateMe(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})
}

// RegisterPaymentServiceClient registers a Payment service client to http server
func RegisterPaymentServiceClient(mux *http.ServeMux, s acourse.PaymentServiceClient) {
	sv := "/acourse.PaymentService"

	mux.HandleFunc(sv+"/ListWaitingPayments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.ListRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListWaitingPayments(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/ListHistoryPayments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.ListRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListHistoryPayments(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/ApprovePayments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.PaymentIDsRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ApprovePayments(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/RejectPayments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.PaymentIDsRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.RejectPayments(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/UpdatePrice", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.PaymentUpdatePriceRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.UpdatePrice(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})
}

// RegisterCourseServiceClient registers a Course service client to http server
func RegisterCourseServiceClient(mux *http.ServeMux, s acourse.CourseServiceClient) {
	sv := "/acourse.CourseService"

	mux.HandleFunc(sv+"/ListCourses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.ListRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListCourses(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/ListPublicCourses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.ListRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListPublicCourses(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/ListOwnCourses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.UserIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListOwnCourses(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/ListEnrolledCourses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.UserIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListEnrolledCourses(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/GetCourse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.CourseIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetCourse(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/CreateCourse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.Course)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.CreateCourse(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/UpdateCourse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.Course)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.UpdateCourse(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/EnrollCourse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.EnrollRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.EnrollCourse(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/AttendCourse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.CourseIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.AttendCourse(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/OpenAttend", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.CourseIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.OpenAttend(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/CloseAttend", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.CourseIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.CloseAttend(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})
}

// RegisterAssignmentServiceClient registers s Assignment service client to http server
func RegisterAssignmentServiceClient(mux *http.ServeMux, s acourse.AssignmentServiceClient) {
	sv := "/acourse.AssignmentService"

	mux.HandleFunc(sv+"/CreateAssignment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.Assignment)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.CreateAssignment(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/UpdateAssignment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.Assignment)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.UpdateAssignment(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/OpenAssignment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.AssignmentIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.OpenAssignment(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/CloseAssignment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.AssignmentIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.CloseAssignment(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/ListAssignments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.CourseIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListAssignments(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/DeleteAssignment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.AssignmentIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.DeleteAssignment(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/SubmitUserAssignment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.UserAssignment)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.SubmitUserAssignment(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/DeleteUserAssignment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.UserAssignmentIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.DeleteUserAssignment(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/GetUserAssignments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.AssignmentIDsRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetUserAssignments(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})

	mux.HandleFunc(sv+"/ListUserAssignments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, httperror.MethodNotAllowed)
			return
		}
		req := new(acourse.CourseIDRequest)
		err := bindJSON(r, req)
		if err != nil {
			handleError(w, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListUserAssignments(makeServiceContext(r), req)
		if err != nil {
			handleError(w, httperror.GRPC(err))
			return
		}
		handleOK(w, res)
	})
}
