package webhandler

import (
	"os"
	"testing"

	. "github.com/akfaew/aeutils"
)

func TestMain(m *testing.M) {
	AppEngineSetup()
	result := m.Run()
	AppEngineShutdown()

	os.Exit(result)
}
