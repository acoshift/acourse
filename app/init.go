package app

import (
	"net/http"
	"strings"

	"errors"

	"github.com/acoshift/go-firebase-admin"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/google/uuid"
)

const (
	keyCurrentUserID = "CurrentUserID"
	keyRequestID     = "RequestID"
)

var (
	firApp  *admin.FirebaseApp
	firAuth *admin.FirebaseAuth

	tokenError = CreateErrors(http.StatusUnauthorized, "token")
)

// InitService inits service
func InitService(service *gin.Engine, projectID string) (err error) {
	firApp, err = admin.InitializeApp(admin.ProjectID(projectID))
	if err != nil {
		return
	}
	firAuth = firApp.Auth()

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
