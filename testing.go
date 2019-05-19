package webhandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akfaew/test"
	"github.com/gorilla/mux"
)

type Response struct {
	t        *testing.T
	Response *httptest.ResponseRecorder
}

func HTTPGet(router *mux.Router, t *testing.T, req *http.Request) *Response {
	t.Helper()

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	return &Response{
		t:        t,
		Response: res,
	}
}

func (r *Response) Fixture() {
	r.t.Helper()

	r.FixtureExtra("")
}

func (r *Response) FixtureExtra(extra string) {
	r.t.Helper()

	body, err := ioutil.ReadAll(r.Response.Body)
	test.NoError(r.t, err)

	code := fmt.Sprintf("Code: %d\n\n", r.Response.Code)
	test.FixtureExtra(r.t, extra, code+string(body))
}

func (r *Response) Status(want int) {
	r.t.Helper()

	if r.Response.Code != want {
		r.t.Fatalf("Status Code == %d (expected %d)", r.Response.Code, want)
	}
}
