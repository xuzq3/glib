package logx

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

	logger2 := NewLogrusLogger(NewOption().SetTextFormat().AddOutput(os.Stdout))
	logger2.Info("hello")
	logger2.WithFields(Fields{
		"a": "b",
		"c": "d",
	}).Infof("hello %s", "world")
}

func BenchmarkLogrusLogger(b *testing.B) {
	logger := NewLogrusLogger(NewOption().SetJsonFormat().AddOutput(os.Stdout))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(Fields{
			"a": "b",
			"c": "d",
		}).Infof("hello %s", "world")
	}
}
