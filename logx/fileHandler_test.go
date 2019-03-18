package logx_test

import (
	"libcore/logx"
	"testing"
)

func TestFileHandler(t *testing.T) {
	fh, err := logx.NewFileHandler("./tmp.log")
	if err != nil {
		t.Log("NewFileHandler error: ", err.Error())
		return
	}
	fh.SetCallDepth(2)
	fh.SetLevel(logx.DEBUG)
	fh.Trace("Trace")
	fh.Tracef("%s", "Tracef")
	fh.Debug("Debug")
	fh.Debugf("%s", "Debugf")
	fh.Info("Info")
	fh.Infof("%s", "Infof")
	fh.Error("Error")
	fh.Errorf("%s", "Errorf")
	fh.Fatal("Fatal")
	fh.Fatalf("%s", "Fatalf")
}
