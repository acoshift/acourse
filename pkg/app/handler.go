package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/acoshift/httperror"
	"github.com/unrolled/render"
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

func handleJSON(w http.ResponseWriter, status int, v interface{}) {
	if err := rr.JSON(w, status, v); err != nil {
		panic(err)
	}
}

func bindJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func handleOK(w http.ResponseWriter, v interface{}) {
	handleJSON(w, http.StatusOK, v)
}

func handleError(w http.ResponseWriter, r error) {
	if err, ok := r.(*httperror.Error); ok {
		log.Println(r)
		handleJSON(w, err.Status, &ErrorReply{err})
	} else {
		handleError(w, httperror.InternalServerErrorWith(r))
	}
}

func handleSuccess(w http.ResponseWriter) {
	handleOK(w, success)
}

func handleHTML(w http.ResponseWriter, name string, binding interface{}) {
	if err := rr.HTML(w, http.StatusOK, name, binding); err != nil {
		panic(err)
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
}
