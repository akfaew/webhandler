// Package webhandler provides utilities for serving HTML pages with AppEngine and html/template
package webhandler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	. "github.com/akfaew/aeutils"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

const (
	ErrorInternal = "Internal server error"
)

var (
	TemplateDir   = "templates/"
	TemplateBase  = "base.html"
	TemplateError = "error.html"
)

var ErrorContextFunc func(message string) interface{} = nil

type WebHandler func(http.ResponseWriter, *http.Request) *WebError

type WebTemplate struct {
	t *template.Template
}

type Templates map[string]*WebTemplate

type WebError struct {
	Code    int
	Error   error  // for the logs
	Message string // for the user
}

func WebErrorf(code int, err error, format string, v ...interface{}) *WebError {
	return &WebError{
		Code:    code,
		Error:   err,
		Message: fmt.Sprintf(format, v...),
	}
}

// ParseTemplate parses the template nesting it in the base template
func ParseTemplate(filename string) *WebTemplate {
	tmpl := template.Must(template.ParseFiles(TemplateDir+TemplateBase, TemplateDir+filename))

	return &WebTemplate{tmpl.Lookup(TemplateBase)}
}

// ServeHTTP renders an error template and logs the failure if a problem occurs.
func (fn WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if e := fn(w, r); e != nil { // e is *WebError, not os.Error.
		// Log for the site admin
		switch e.Code / 100 {
		case 4:
			LogInfof(ctx, "Handler error %d. err=\"%v\", msg=\"%s\"",
				e.Code, e.Error, e.Message)
		default:
			LogErrorf(ctx, "Handler error %d. err=\"%v\", msg=\"%s\"",
				e.Code, e.Error, e.Message)
		}

		// Error for the user
		tmpl := template.Must(template.ParseFiles(TemplateDir + TemplateError))
		buf := new(bytes.Buffer)
		var vars interface{}
		if ErrorContextFunc != nil {
			vars = ErrorContextFunc(e.Message)
		}
		if err := tmpl.Execute(buf, vars); err != nil {
			LogErrorfd(ctx, "err=%v", err)
			http.Error(w, "Internal error executing template", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(e.Code)
		if _, err := buf.WriteTo(w); err != nil {
			LogErrorfd(ctx, "err=%v", err)
			http.Error(w, "Internal error executing template", http.StatusInternalServerError)
			return
		}
	}
}

func WebHandle(r *mux.Router, method string, path string, handler func(w http.ResponseWriter, r *http.Request) *WebError) {
	r.Methods(method).Path(path).Handler(WebHandler(handler))
}

// Executor uses webContext() to obtain a list of variables to pass to the underlying html template.
func (tmpl *WebTemplate) Executor(webContext func(*http.Request, *WebTemplate) (interface{}, *WebError)) func(http.ResponseWriter, *http.Request) *WebError {
	return func(w http.ResponseWriter, r *http.Request) *WebError {
		// Get the web context
		vars, weberr := webContext(r, tmpl)
		if weberr != nil {
			return weberr
		}

		// Execute the template
		buf := new(bytes.Buffer)
		if err := tmpl.t.Execute(buf, vars); err != nil {
			return WebErrorf(http.StatusInternalServerError, err, "Internal error executing template")
		}

		// Return the result
		if _, err := buf.WriteTo(w); err != nil {
			return WebErrorf(http.StatusInternalServerError, err, "")
		}

		return nil
	}
}
