package logger

import (
	"os"
	"testing"
)

func TestLogrusLogger(t *testing.T) {
	logger := NewLogrusLogger(NewOption().SetJsonFormat().AddOutput(os.Stdout))
	logger.Info("hello")
	logger.WithFields(Fields{
		"a": "b",
		"c": "d",
	}).Infof("hello %s", "world")
}

func BenchmarkLogrusLogger(b *testing.B) {
	logger := NewLogrusLogger(NewOption().SetJsonFormat().AddOutput(os.Stdout))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//logger.Info("hello world")
		logger.WithFields(Fields{
			"a": "b",
			"c": "d",
		}).Infof("hello %s", "world")
	}
}