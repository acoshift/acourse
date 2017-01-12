package app

import (
	"github.com/acoshift/httperror"
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

// RegisterUserService registers a User service
func RegisterUserService(service *gin.Engine, s UserService) {
	service.POST("/acourse.UserService/GetUsers", func(ctx *gin.Context) {
		req := new(IDsRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		res, err := s.GetUsers(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res.Expose())
	})

	service.GET("/acourse.UserService/GetMe", func(ctx *gin.Context) {
		res, err := s.GetMe(ctx.Request.Context())
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res.Expose())
	})

	service.POST("/acourse.UserService/UpdateMe", func(ctx *gin.Context) {
		req := new(UserRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, httperror.BadRequestWith(err))
			return
		}
		err = s.UpdateMe(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
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
