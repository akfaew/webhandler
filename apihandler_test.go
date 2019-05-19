package webhandler

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/akfaew/aeutils"
	"github.com/akfaew/test"
	"github.com/gorilla/mux"
)

func apiSuccess(w http.ResponseWriter, r *http.Request) *APIError {
	return nil
}

func apiFailure(w http.ResponseWriter, r *http.Request) *APIError {
	return NewAPIError(http.StatusInternalServerError, fmt.Errorf("ooups"))
}

func apiRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	APIHandle(r, http.MethodGet, "/success", apiSuccess)
	APIHandle(r, http.MethodGet, "/failure", apiFailure)

	return r
}

func TestAPIHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req, err := Inst.NewRequest(http.MethodGet, "/success", nil)
		test.NoError(t, err)

		HTTPGetRouter(apiRouter(), t, req).Fixture()
	})

	t.Run("Failure", func(t *testing.T) {
		req, err := Inst.NewRequest(http.MethodGet, "/failure", nil)
		test.NoError(t, err)

		HTTPGetRouter(apiRouter(), t, req).Fixture()
	})
}
