package webhandler

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/akfaew/aeutils"
	"github.com/akfaew/test"
	"github.com/gorilla/mux"
)

var simple = ParseTemplate("simple.html")
var failure = ParseTemplate("simple.html")

func webContext(r *http.Request, tmpl *WebTemplate) (interface{}, *WebError) {
	if tmpl == failure {
		return nil, WebErrorf(http.StatusInternalServerError, fmt.Errorf("ooups"), "User error")
	}

	return struct {
		Title string
	}{
		Title: "Template Title",
	}, nil
}

func errorContext(message string) interface{} {
	return struct {
		Message string
	}{
		Message: message,
	}
}

func webRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	ErrorContextFunc = errorContext

	WebHandle(r, http.MethodGet, "/success", simple.Executor(webContext))
	WebHandle(r, http.MethodGet, "/failure", failure.Executor(webContext))

	return r
}

func TestWebHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req, err := Inst.NewRequest(http.MethodGet, "/success", nil)
		test.NoError(t, err)

		HTTPGetRouter(t, webRouter(), req).Fixture()
	})

	t.Run("Failure", func(t *testing.T) {
		req, err := Inst.NewRequest(http.MethodGet, "/failure", nil)
		test.NoError(t, err)

		HTTPGetRouter(t, webRouter(), req).Fixture()
	})
}
