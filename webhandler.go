// Package webhandler provides utilities for serving HTML pages with AppEngine and html/template.
//
// Additionally, utilities for communicating with APIs are provided (no user error reporting).
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
	Code    int    // HTTP response code
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
func ParseTemplate(filenames ...string) *WebTemplate {
	if len(filenames) == 0 {
		return nil
	}

	paths := []string{}
	for _, f := range filenames {
		paths = append(paths, TemplateDir+f)
	}
	tmpl := template.Must(template.ParseFiles(paths...))

	return &WebTemplate{tmpl.Lookup(filenames[0])}
}

// ServeHTTP renders an error template and logs the failure if a problem occurs.
func (fn WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *WebError, not os.Error.
		ctx := appengine.NewContext(r)

		// Log for the site admin. No logging occurs if no error is passed.
		if e.Error != nil || len(e.Message) > 0 {
			switch e.Code / 100 {
			case 4:
				LogInfof(ctx, "Handler error %d. err=\"%v\", msg=\"%s\"",
					e.Code, e.Error, e.Message)
			default:
				LogErrorf(ctx, "Handler error %d. err=\"%v\", msg=\"%s\"",
					e.Code, e.Error, e.Message)
			}
		}

		// Error for the user
		tmpl, err := template.ParseFiles(TemplateDir + TemplateError)
		if err != nil { // if TemplateError does not exist
			http.Error(w, e.Message, e.Code)
			return
		}

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
func (tmpl *WebTemplate) Executor(webContext func(http.ResponseWriter, *http.Request, *WebTemplate) (interface{}, *WebError)) func(http.ResponseWriter, *http.Request) *WebError {
	return func(w http.ResponseWriter, r *http.Request) *WebError {
		// Get the web context
		vars, weberr := webContext(w, r, tmpl)
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
