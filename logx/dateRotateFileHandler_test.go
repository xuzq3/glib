package logx_test

import (
	"libcore/logx"
	"testing"
	"time"
)

func TestDateRotateFileHandler(t *testing.T) {
	fh, err := logx.NewDateRotateFileHandler("./", "tmp.log", 4)
	if err != nil {
		t.Log("NewDateRotateFileHandler error: ", err.Error())
		return
	}
	fh.SetCallDepth(3)
	fh.SetLevel(logx.DEBUG)
	for i := 0; i < 10; i++ {
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
		time.Sleep(time.Second * 1)
	}
}
