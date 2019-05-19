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

type WebTestResponse struct {
	t                *testing.T
	ResponseRecorder *httptest.ResponseRecorder
}

// Use like this:
//
// func HTTPGet(t *testing.T, req *http.Request) *WebTestResponse {
// 	return webhandler.HTTPGetRouter(t, Router(), req)
// }
func HTTPGetRouter(t *testing.T, router *mux.Router, req *http.Request) *WebTestResponse {
	t.Helper()

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	return &WebTestResponse{
		t:                t,
		ResponseRecorder: res,
	}
}

func (r *WebTestResponse) Fixture() {
	r.t.Helper()

	r.FixtureExtra("")
}

func (r *WebTestResponse) FixtureExtra(extra string) {
	r.t.Helper()

	body, err := ioutil.ReadAll(r.ResponseRecorder.Body)
	test.NoError(r.t, err)

	code := fmt.Sprintf("Code: %d\n\n", r.ResponseRecorder.Code)
	test.FixtureExtra(r.t, extra, code+string(body))
}

func (r *WebTestResponse) Status(want int) {
	r.t.Helper()

	if r.ResponseRecorder.Code != want {
		r.t.Fatalf("Status Code == %d (expected %d)", r.ResponseRecorder.Code, want)
	}
}
