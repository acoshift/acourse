package app

import (
	"github.com/acoshift/httperror"
	"gopkg.in/gin-gonic/gin.v1"
)

// MountHealthService mounts a Health service
func MountHealthService(service *gin.Engine, s HealthService) {
	service.GET("/acourse.HealthService/Check", func(ctx *gin.Context) {
		err := s.Check(ctx.Request.Context())
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleSuccess(ctx)
	})
}

// MountUserService mounts a User service
func MountUserService(service *gin.Engine, s UserService) {
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
		handleOK(ctx, res)
	})

	service.GET("/acourse.UserService/GetMe", func(ctx *gin.Context) {
		r, err := s.GetMe(ctx.Request.Context())
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, r)
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

// MountPaymentService mounts a payment service
func MountPaymentService(service *gin.Engine, s PaymentService) {
	service.POST("/acourse.PaymentService/ListPayments", func(ctx *gin.Context) {
		req := new(PaymentListRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		res, err := s.ListPayments(ctx.Request.Context(), req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		handleOK(ctx, res)
	})

	service.POST("/acourse.PaymentService/ApprovePayments", func(ctx *gin.Context) {
		req := new(IDsRequest)
		err := ctx.BindJSON(req)
		if err != nil {
			handleError(ctx, err)
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
			handleError(ctx, err)
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
