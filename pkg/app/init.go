package app

import (
	"net/http"
	"strings"

	"errors"

	"github.com/acoshift/go-firebase-admin"
	"github.com/google/uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

const (
	keyCurrentUserID = "CurrentUserID"
	keyRequestID     = "RequestID"
)

var (
	firAuth *admin.FirebaseAuth

	tokenError = CreateErrors(http.StatusUnauthorized, "token")
)

// InitService inits service
func InitService(service *gin.Engine, auth *admin.FirebaseAuth) (err error) {
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
		return
	}
	claims, err := firAuth.VerifyIDToken(tk[1])
	if err != nil {
		handleError(ctx, tokenError(err))
		return
	}
	ctx.Set(keyCurrentUserID, claims.UserID)
	ctx.Next()
}

func requestIDMiddleware(ctx *gin.Context) {
	rid := uuid.New().String()
	ctx.Header("X-Request-Id", rid)
	ctx.Set(keyRequestID, rid)
	ctx.Next()
}
