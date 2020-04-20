package logger

import "testing"

func TestFacade(t *testing.T) {
	cfg := Config{
		Level: InfoLevel,
		JsonFormat: true,
	}
	err := Init(cfg)
	if err != nil {
		t.Error(err)
	}
	Info("hello")
	WithFields(Fields{
		"a": 1,
		"c": "C",
	}).Info("info")
}
