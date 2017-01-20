package app

import (
	"net/http"
	"os"
)

// HealthController is the controller interface for health check
type HealthController interface {
	Check() error
}

// MountHealthController mounts a Health controller to the http server
func MountHealthController(mux *http.ServeMux, c HealthController) {
	mux.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := c.Check()
			if err != nil {
				handleError(w, err)
				return
			}
			handleSuccess(w)
		}
	})
}

type fileFS struct {
	http.FileSystem
}

type noDirFile struct {
	http.File
}

func (fs *fileFS) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return &noDirFile{f}, nil
}

func (f *noDirFile) Readdir(int) ([]os.FileInfo, error) {
	return nil, nil
}

// RenderController is the controller interface for render actions
type RenderController interface {
	Index(*RenderIndexContext) (interface{}, error)
	Course(*RenderCourseContext) (interface{}, error)
}

// MountRenderController mount a Render template controller on the given resource
func MountRenderController(mux *http.ServeMux, c RenderController) {
	// cc := func(ctx *gin.Context) {
	// 	ctx.Header("Cache-Control", "public, max-age=31536000")
	// }

	mux.Handle("/static/*filepath", http.StripPrefix("/static", http.FileServer(&fileFS{http.Dir("public")})))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/acourse-120.png")
	})

	// server.Group("/static", cc).Static("", "public")

	// server.StaticFile("/favicon.ico", "public/acourse-120.png")

	// server.GET("/course/:courseID", func(ctx *gin.Context) {
	// 	rctx, err := NewRenderCourseContext(ctx)
	// 	if err != nil {
	// 		handleError(ctx, err)
	// 		return
	// 	}
	// 	res, err := c.Course(rctx)
	// 	if err != nil {
	// 		handleError(ctx, err)
	// 		return
	// 	}
	// 	if res == nil {
	// 		handleRedirect(ctx, "/")
	// 		return
	// 	}
	// 	handleHTML(ctx, "index", res)
	// })

	// h := func(ctx *gin.Context) {
	// 	rctx, err := NewRenderIndexContext(ctx)
	// 	if err != nil {
	// 		handleError(ctx, err)
	// 		return
	// 	}
	// 	res, err := c.Index(rctx)
	// 	if err != nil {
	// 		handleError(ctx, err)
	// 		return
	// 	}
	// 	handleHTML(ctx, "index", res)
	// }

	// server.Use(h)
}
