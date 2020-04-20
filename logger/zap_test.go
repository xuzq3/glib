package logger

import (
	"os"
	"testing"
)

func TestZapLogger(t *testing.T) {
	logger := NewZapLogger(NewOption().SetJsonFormat().AddOutput(os.Stdout))
	logger.Infof("hello")
	logger.WithFields(Fields{
		"a": "b",
		"c": "d",
	}).Infof("hello %s", "world")
}

func BenchmarkZapLogger(b *testing.B) {
	logger := NewZapLogger(NewOption().SetJsonFormat().AddOutput(os.Stdout))
	for i := 0; i < b.N; i++ {
		//logger.Infof("hello")
		logger.WithFields(Fields{
			"a": "b",
			"c": "d",
		}).Infof("hello %s", "world")
	}
}