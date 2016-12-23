package app

import (
	"net/http"
	"strings"

	"github.com/acoshift/go-firebase-admin"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

const (
	keyCurrentUserID = "CurrentUserID"
)

var (
	firApp  *admin.FirebaseApp
	firAuth *admin.FirebaseAuth

	tokenError = CreateErrors(http.StatusUnauthorized, "token")
)

// InitService inits service
func InitService(service *echo.Echo, projectID string) (err error) {
	firApp, err = admin.InitializeApp(admin.ProjectID(projectID))
	if err != nil {
		return
	}
	firAuth = firApp.Auth()

	service.Use(requestIDMiddleware)
	service.Use(jwtMiddleware)
	return
}

func jwtMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		auth := strings.TrimSpace(ctx.Request().Header.Get("Authorization"))
		if len(auth) == 0 {
			return h(ctx)
		}
		tk := strings.Split(auth, " ")
		if len(tk) != 2 || strings.ToLower(tk[0]) != "bearer" {
			return handleError(ctx, tokenError("invalid authorization header"))
		}
		claims, err := firAuth.VerifyIDToken(tk[1])
		if err != nil {
			return handleError(ctx, tokenError(err))
		}
		ctx.Set(keyCurrentUserID, claims.UserID)
		return h(ctx)
	}
}

func requestIDMiddleware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("X-Request-Id", uuid.New().String())
		return h(ctx)
	}
}
