package logx_test

import (
	"libcore/logx"
	"testing"
)

func TestConsoleHandler(t *testing.T) {
	ch := logx.NewConsoleHandler()
	ch.SetCallDepth(3)
	ch.SetLevel(logx.DEBUG)
	ch.Trace("Trace")
	ch.Debug("Debug")
	ch.Debugf("%s", "Debugf")
	ch.Info("TraInfoce")
	ch.Infof("%s", "Infof")
	ch.Error("Error")
	ch.Errorf("%s", "Errorf")
	ch.Fatal("Fatal")
	ch.Fatalf("%s", "Fatalf")
}
