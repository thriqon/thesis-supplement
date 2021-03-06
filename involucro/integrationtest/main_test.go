package integrationtest

import (
	"os"
	"testing"

	"github.com/thriqon/involucro/ilog"
)

func TestMain(m *testing.M) {
	if !testing.Verbose() {
		ilog.StdLog.SetPrintFunc(func(b ilog.Bough) {})
	}
	os.Exit(m.Run())
}
