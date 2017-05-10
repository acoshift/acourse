package view

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/acoshift/header"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

const templateDir = "template"

var (
	m         = minify.New()
	muExecute = &sync.Mutex{}
	templates = make(map[interface{}]*templateStruct)
)

type templateStruct struct {
	*template.Template
	set []string
}

func init() {
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)

	parseTemplate(keyIndex, []string{"index.tmpl", "layout.tmpl"})
}

func joinTemplateDir(files []string) []string {
	r := make([]string, len(files))
	for i, f := range files {
		r[i] = filepath.Join(templateDir, f)
	}
	return r
}

func parseTemplate(key interface{}, set []string) {
	templateName := strings.TrimSuffix(set[0], ".tmpl")
	t := template.New("")
	t.Funcs(template.FuncMap{
		"templateName": func() string {
			return templateName
		},
	})
	_, err := t.ParseFiles(joinTemplateDir(set)...)
	if err != nil {
		log.Fatalf("internal: parse template %s error; %v", templateName, err)
	}
	t = t.Lookup("root")
	if t == nil {
		log.Fatalf("internal: root template not found in %s", templateName)
	}
	templates[key] = &templateStruct{
		Template: t,
		set:      set,
	}
}

func render(w http.ResponseWriter, r *http.Request, key, data interface{}) {
	t := templates[key]
	if t == nil {
		http.Error(w, fmt.Sprintf("template not found"), http.StatusInternalServerError)
		return
	}
	if dev {
		muExecute.Lock()
		defer muExecute.Unlock()
		parseTemplate(key, t.set)
		t = templates[key]
	}

	w.Header().Set(header.ContentType, "text/html; charset=utf-8")
	pipe := &bytes.Buffer{}
	err := t.Execute(pipe, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = m.Minify("text/html", w, pipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
