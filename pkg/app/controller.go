package app

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// HealthController is the controller interface for health check
type HealthController interface {
	Check() error
}

// MountHealthController mounts a Health controller to the http server
func MountHealthController(mux *http.ServeMux, c HealthController) {
	mux.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		err := c.Check()
		if err != nil {
			handleError(w, err)
			return
		}
		handleSuccess(w)
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
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/acourse-120.png")
	})

	pushList := []string{}
	{
		dirs := []string{"js", "css"}
		for _, d := range dirs {
			fl, _ := ioutil.ReadDir("public/" + d)
			for _, f := range fl {
				pushList = append(pushList, "/static/"+d+"/"+f.Name())
			}
		}
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static") {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			http.StripPrefix("/static", http.FileServer(&fileFS{http.Dir("public")})).ServeHTTP(w, r)
			return
		}

		if p, ok := w.(http.Pusher); ok {
			opt := &http.PushOptions{
				Method: http.MethodGet,
				Header: r.Header,
			}
			for _, t := range pushList {
				p.Push(t, opt)
			}
		}

		if strings.HasPrefix(r.URL.Path, "/course/") && !strings.Contains(r.URL.Path[8:], "/") {
			http.StripPrefix("/course/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx, err := NewRenderCourseContext(r)
				if err != nil {
					handleError(w, err)
					return
				}
				res, err := c.Course(ctx)
				if err != nil {
					handleError(w, err)
					return
				}
				if res == nil {
					handleRedirect(w, r, "/")
				}
				handleHTML(w, "index", res)
			})).ServeHTTP(w, r)
			return
		}

		ctx, err := NewRenderIndexContext(r)
		if err != nil {
			handleError(w, err)
			return
		}
		res, err := c.Index(ctx)
		if err != nil {
			handleError(w, err)
			return
		}
		handleHTML(w, "index", res)
	})
}
