package webhandler

import (
	"net/http"

	. "github.com/akfaew/aeutils"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

type APIHandler func(http.ResponseWriter, *http.Request) *APIError

type APIError struct {
	Code  int   // HTTP response code
	Error error // for the logs
}

func NewAPIError(code int, error error) *APIError {
	return &APIError{
		Code:  code,
		Error: error,
	}
}

func (fn APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		ctx := appengine.NewContext(r)

		LogErrorf(ctx, "Handler returned code %d with error %v.", e.Code, e.Error)
		http.Error(w, "", e.Code)
	}
}

func APIHandle(r *mux.Router, method string, path string, handler func(w http.ResponseWriter, r *http.Request) *APIError) {
	r.Methods(method).Path(path).Handler(APIHandler(handler))
}
