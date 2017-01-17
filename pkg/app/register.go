package app

import (
	"context"
	"net/http"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/httperror"
	"google.golang.org/grpc/metadata"
	"gopkg.in/gin-gonic/gin.v1"
)

func makeServiceContext(r *http.Request) context.Context {
	md := metadata.MD{}
	header := r.Header
	if v := header.Get("Authorization"); v != "" {
		md = metadata.Join(md, metadata.Pairs("authorization", v))
	}
	return metadata.NewContext(context.Background(), md)
}

// RegisterUserServiceClient registers a User service client to http server
func RegisterUserServiceClient(httpServer *gin.Engine, s acourse.UserServiceClient) {
	sv := "/acourse.UserService"

	httpServer.POST(sv+"/GetUser", func(ctx *gin.Context) {
		req := new(acourse.GetUserRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetUser(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/GetMe", func(ctx *gin.Context) {
		res, err := s.GetMe(makeServiceContext(ctx.Request), new(acourse.Empty))
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/UpdateMe", func(ctx *gin.Context) {
		req := new(acourse.User)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.UpdateMe(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
}

// RegisterEmailServiceClient registers a Email service client to http server
func RegisterEmailServiceClient(httpServer *gin.Engine, s acourse.EmailServiceClient) {
	sv := "/acourse.EmailService"

	httpServer.POST(sv+"/Send", func(ctx *gin.Context) {
		req := new(acourse.Email)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.Send(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
}

// RegisterPaymentServiceClient registers a Payment service client to http server
func RegisterPaymentServiceClient(httpServer *gin.Engine, s acourse.PaymentServiceClient) {
	sv := "/acourse.PaymentService"

	httpServer.POST(sv+"/ListWaitingPayments", func(ctx *gin.Context) {
		req := new(acourse.ListRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListWaitingPayments(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/ListHistoryPayments", func(ctx *gin.Context) {
		req := new(acourse.ListRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListHistoryPayments(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/ApprovePayments", func(ctx *gin.Context) {
		req := new(acourse.PaymentIDsRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.ApprovePayments(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})

	httpServer.POST(sv+"/RejectPayments", func(ctx *gin.Context) {
		req := new(acourse.PaymentIDsRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.RejectPayments(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})

	httpServer.POST(sv+"/UpdatePrice", func(ctx *gin.Context) {
		req := new(acourse.PaymentUpdatePriceRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.UpdatePrice(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
}

// RegisterCourseServiceClient registers a Course service client to http server
func RegisterCourseServiceClient(httpServer *gin.Engine, s acourse.CourseServiceClient) {
	sv := "/acourse.CourseService"

	httpServer.POST(sv+"/ListCourses", func(ctx *gin.Context) {
		req := new(acourse.ListRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListCourses(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/ListPublicCourses", func(ctx *gin.Context) {
		req := new(acourse.ListRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListPublicCourses(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/ListOwnCourses", func(ctx *gin.Context) {
		req := new(acourse.UserIDRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListOwnCourses(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/ListEnrolledCourses", func(ctx *gin.Context) {
		req := new(acourse.UserIDRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListEnrolledCourses(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/GetCourse", func(ctx *gin.Context) {
		req := new(acourse.CourseIDRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetCourse(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/CreateCourse", func(ctx *gin.Context) {
		req := new(acourse.Course)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.CreateCourse(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST(sv+"/UpdateCourse", func(ctx *gin.Context) {
		req := new(acourse.Course)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.UpdateCourse(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})

	httpServer.POST(sv+"/EnrollCourse", func(ctx *gin.Context) {
		req := new(acourse.EnrollRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.EnrollCourse(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})

	httpServer.POST(sv+"/AttendCourse", func(ctx *gin.Context) {
		req := new(acourse.CourseIDRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.AttendCourse(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})

	httpServer.POST(sv+"/OpenAttend", func(ctx *gin.Context) {
		req := new(acourse.CourseIDRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.OpenAttend(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})

	httpServer.POST(sv+"/CloseAttend", func(ctx *gin.Context) {
		req := new(acourse.CourseIDRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		_, err = s.CloseAttend(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
}

// RegisterAssignmentServiceClient registers s Assignment service client to http server
func RegisterAssignmentServiceClient(httpServer *gin.Engine, s acourse.AssignmentServiceClient) {
	sv := "/acourse.AssignmentService"

	httpServer.POST(sv+"/ListMyAssignmentsByCourse", func(ctx *gin.Context) {
		req := new(acourse.CourseIDRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListMyAssignmentsByCourse(makeServiceContext(ctx.Request), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})
}
