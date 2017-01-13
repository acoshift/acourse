package app

import (
	"log"

	"github.com/acoshift/acourse/pkg/acourse"
	"github.com/acoshift/httperror"
	"google.golang.org/grpc"
	"gopkg.in/gin-gonic/gin.v1"
)

// RegisterHealthService registers a Health service
func RegisterHealthService(service *gin.Engine, s HealthService) {
	service.GET("/acourse.HealthService/Check", func(ctx *gin.Context) {
		err := s.Check(ctx.Request.Context())
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
}

// RegisterUserServiceServer registers a User service server
func RegisterUserServiceServer(httpServer *gin.Engine, addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	s := acourse.NewUserServiceClient(conn)
	httpServer.POST("/acourse.UserService/GetUser", func(ctx *gin.Context) {
		req := new(acourse.GetUserRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetUser(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	// httpServer.GET("/acourse.UserService/GetMe", func(ctx *gin.Context) {
	// 	res, err := s.GetMe(ctx.Request.Context())
	// 	if err != nil {
	// 		handleError(ctx, err)
	// 		return
	// 	}
	// 	handleOK(ctx, res)
	// })

	// httpServer.POST("/acourse.UserService/UpdateMe", func(ctx *gin.Context) {
	// 	req := new(UserRequest)
	// 	err := ctx.BindJSON(req)
	// 	if err != nil {
	// 		handleError(ctx, httperror.BadRequestWith(err))
	// 		return
	// 	}
	// 	err = s.UpdateMe(ctx.Request.Context(), req)
	// 	if err != nil {
	// 		handleError(ctx, err)
	// 		return
	// 	}
	// 	handleSuccess(ctx)
	// })
}

// RegisterPaymentService registers a payment service
func RegisterPaymentService(service *gin.Engine, s PaymentService) {
	service.POST("/acourse.PaymentService/ListPayments", func(ctx *gin.Context) {
		req := new(PaymentListRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListPayments(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res.Expose())
	})

	service.POST("/acourse.PaymentService/ApprovePayments", func(ctx *gin.Context) {
		req := new(IDsRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		err = s.ApprovePayments(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})

	service.POST("/acourse.PaymentService/RejectPayments", func(ctx *gin.Context) {
		req := new(IDsRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		err = s.RejectPayments(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
}

// RegisterCourseService registers a course service
func RegisterCourseService(service *gin.Engine, s CourseService) {
	service.POST("/acourse.CourseService/ListCourses", func(ctx *gin.Context) {
		req := new(CourseListRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.ListCourses(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res.Expose())
	})

	service.GET("/acourse.CourseService/ListEnrolledCourses", func(ctx *gin.Context) {
		res, err := s.ListEnrolledCourses(ctx.Request.Context())
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res.Expose())
	})
}
