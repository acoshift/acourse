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
	if err := rr.JSON(ctx.Writer, http.StatusOK, r); err != nil {
		panic(err)
	}
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
	handleOK(ctx, success)
}

func handleHTML(ctx *gin.Context, name string, binding interface{}) {
	if err := rr.HTML(ctx.Writer, http.StatusOK, name, binding); err != nil {
		panic(err)
	}
}

func handleRedirect(ctx *gin.Context, path string) {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusFound)
}
