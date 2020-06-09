package logx

import (
	"os"
	"testing"
)

func TestZapLogger(t *testing.T) {
	logger := NewZapLogger(NewOption().SetTextFormat().AddOutput(os.Stdout))
	logger.Infof("hello")
	logger.WithFields(Fields{
		"a": "b",
		"c": "d",
	}).Infof("hello %s", "world")

	logger2 := NewZapLogger(NewOption().SetTextFormat().AddOutput(os.Stdout))
	logger2.Info("hello")
	logger2.WithFields(Fields{
		"a": "b",
		"c": "d",
	}).Infof("hello %s", "world")
}

func BenchmarkZapLogger(b *testing.B) {
	logger := NewZapLogger(NewOption().SetJsonFormat().AddOutput(os.Stdout))
	for i := 0; i < b.N; i++ {
		logger.WithFields(Fields{
			"a": "b",
			"c": "d",
		}).Infof("hello %s", "world")
	}
}
