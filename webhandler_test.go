package webhandler

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/akfaew/aeutils"
	"github.com/akfaew/test"
	"github.com/gorilla/mux"
)

func webSuccess(w http.ResponseWriter, r *http.Request) *WebError {
	return nil
}

func webFailure(w http.ResponseWriter, r *http.Request) *WebError {
	return NewAPIError(http.StatusInternalServerError, fmt.Errorf("ooups"))
}

func webRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	WebHandle(r, http.MethodGet, "/success", apiSuccess)
	WebHandle(r, http.MethodGet, "/failure", apiFailure)

	return r
}

func TestAPIHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req, err := Inst.NewRequest(http.MethodGet, "/success", nil)
		test.NoError(t, err)

		HTTPGet(APIRouter(), t, req).Fixture()
	})

	t.Run("Failure", func(t *testing.T) {
		req, err := Inst.NewRequest(http.MethodGet, "/failure", nil)
		test.NoError(t, err)

		HTTPGet(APIRouter(), t, req).Fixture()
	})
}
