package app

import (
	"log"
	"net/http"

	"github.com/acoshift/httperror"
	"github.com/unrolled/render"
	"gopkg.in/gin-gonic/gin.v1"
)

// ErrorReply is the error response
type ErrorReply struct {
	Error *httperror.Error `json:"error"`
}

// SuccessReply is the success response without any content
type SuccessReply struct {
	OK int `json:"ok"`
}

var success = &SuccessReply{1}

var rr = render.New(render.Options{DisableHTTPErrorRendering: true})

func handleOK(ctx *gin.Context, r interface{}) {
	ctx.JSON(http.StatusOK, r)
}

func handleError(ctx *gin.Context, r error) {
	if err, ok := r.(*httperror.Error); ok {
		log.Println(r)
		ctx.JSON(err.Status, &ErrorReply{err})
	} else {
		handleError(ctx, httperror.InternalServerErrorWith(r))
	}
}

func handleSuccess(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, success)
}

func handleHTML(ctx *gin.Context, name string, binding interface{}) error {
	return rr.HTML(ctx.Writer, http.StatusOK, name, binding)
}

func handleRedirect(ctx *gin.Context, path string) error {
	ctx.Redirect(http.StatusTemporaryRedirect, path)
	return nil
}
