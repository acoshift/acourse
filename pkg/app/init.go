package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/acoshift/go-firebase-admin"
	"github.com/acoshift/httperror"
	"github.com/google/uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

// ContextKey is the key for app's context
type ContextKey int

// Predefined context keys
const (
	KeyCurrentUserID ContextKey = iota
	KeyRequestID
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

func validateHeaderToken(header string) (string, error) {
	tk := strings.Split(header, " ")
	if len(tk) != 2 || strings.ToLower(tk[0]) != "bearer" {
		return "", errors.New("invalid authorization header")
	}
	claims, err := firAuth.VerifyIDToken(tk[1])
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

func jwtMiddleware(ctx *gin.Context) {
	auth := strings.TrimSpace(ctx.Request.Header.Get("Authorization"))
	if len(auth) == 0 {
		ctx.Next()
		return
	}
	userID, err := validateHeaderToken(auth)
	if err != nil {
		handleError(ctx, tokenError(err))
		ctx.Abort()
		return
	}
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), KeyCurrentUserID, userID))
	ctx.Set(keyCurrentUserID, userID)
}

func requestIDMiddleware(ctx *gin.Context) {
	rid := uuid.New().String()
	ctx.Header("X-Request-Id", rid)
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), KeyRequestID, rid))
	ctx.Set(keyRequestID, rid)
}
