package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/httperror"
	"github.com/google/uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

const (
	keyCurrentUserID = "CurrentUserID"
	keyRequestID     = "RequestID"
)

var (
	firAuth *admin.Auth

	tokenError = httperror.New(http.StatusUnauthorized, "token")
)

// InitService inits service
func InitService(service *gin.Engine, auth *admin.Auth) (err error) {
	firAuth = auth

	service.Use(requestIDMiddleware)
	service.Use(jwtMiddleware)
	return
}

func jwtMiddleware(ctx *gin.Context) {
	auth := strings.TrimSpace(ctx.Request.Header.Get("Authorization"))
	if len(auth) == 0 {
		ctx.Next()
		return
	}
	tk := strings.Split(auth, " ")
	if len(tk) != 2 || strings.ToLower(tk[0]) != "bearer" {
		handleError(ctx, tokenError(errors.New("invalid authorization header")))
		ctx.Abort()
		return
	}
	claims, err := firAuth.VerifyIDToken(tk[1])
	if err != nil {
		handleError(ctx, tokenError(err))
		ctx.Abort()
		return
	}
	ctx.Set(keyCurrentUserID, claims.UserID)
}

func requestIDMiddleware(ctx *gin.Context) {
	rid := uuid.New().String()
	ctx.Header("X-Request-Id", rid)
	ctx.Set(keyRequestID, rid)
}
