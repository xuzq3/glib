package logx_test

import (
	"libcore/logx"
	"testing"
)

func TestSizeRotateFileHandler(t *testing.T) {
	fh, err := logx.NewSizeRotateFileHandler("./", "tmp.log", 4, 1024*1024)
	if err != nil {
		t.Log("NewSizeRotateFileHandler error: ", err.Error())
		return
	}
	fh.SetCallDepth(2)
	fh.SetLevel(logx.DEBUG)
	for i := 0; i < 1000; i++ {
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
}
