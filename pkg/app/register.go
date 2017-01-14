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
	httpServer.POST("/acourse.UserService/GetUser", func(ctx *gin.Context) {
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

	httpServer.POST("/acourse.UserService/GetMe", func(ctx *gin.Context) {
		res, err := s.GetMe(makeServiceContext(ctx.Request), new(acourse.Empty))
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	httpServer.POST("/acourse.UserService/UpdateMe", func(ctx *gin.Context) {
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
	httpServer.POST("/acourse.EmailService/Send", func(ctx *gin.Context) {
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
	httpServer.POST("/acourse.PaymentService/ListWaitingPayments", func(ctx *gin.Context) {
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

	httpServer.POST("/acourse.PaymentService/ListHistoryPayments", func(ctx *gin.Context) {
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

	httpServer.POST("/acourse.PaymentService/ApprovePayments", func(ctx *gin.Context) {
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

	httpServer.POST("/acourse.PaymentService/RejectPayments", func(ctx *gin.Context) {
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
}

// RegisterCourseServiceClient registers a Course service client to http server
func RegisterCourseServiceClient(service *gin.Engine, s acourse.CourseServiceClient) {
	service.POST("/acourse.CourseService/ListCourses", func(ctx *gin.Context) {
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

	service.POST("/acourse.CourseService/ListPublicCourses", func(ctx *gin.Context) {
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

	service.POST("/acourse.CourseService/ListOwnCourses", func(ctx *gin.Context) {
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

	service.POST("/acourse.CourseService/ListEnrolledCourses", func(ctx *gin.Context) {
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

	service.POST("/acourse.CourseService/GetCourse", func(ctx *gin.Context) {
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
}
