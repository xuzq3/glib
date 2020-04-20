package logger

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	opt    *Option
	logger *zap.SugaredLogger
}

func NewZapLogger(opt *Option) *ZapLogger {
	l := &ZapLogger{
		opt: opt,
	}
	l.initLogger()
	return l
}

func (l *ZapLogger) initLogger() {
	encfg := l.getEncoderConfig()
	var enc zapcore.Encoder
	if l.opt.jsonFormat {
		enc = zapcore.NewJSONEncoder(encfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encfg)
	}

	var level zapcore.Level
	switch strings.ToLower(l.opt.level) {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	default:
		level = zapcore.DebugLevel
	}

	cores := make([]zapcore.Core, 0)
	//if !l.opt.disableConsole {
	//	ws := zapcore.Lock(os.Stdout)
	//	core := zapcore.NewCore(enc, ws, level)
	//	cores = append(cores, core)
	//}
	for _, out := range l.opt.outs {
		ws := zapcore.AddSync(out)
		core := zapcore.NewCore(enc, ws, level)
		cores = append(cores, core)
	}
	combinedCore := zapcore.NewTee(cores...)

	const zapCallerDepth = 1
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(zapCallerDepth + l.opt.callerSkip),
		//zap.Hooks(),
	}
	logger := zap.New(combinedCore, options...).Sugar()
	l.logger = logger
}

func (l *ZapLogger) getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		//EncodeCaller:   l.callerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.000"))
		},
	}
	return encoderConfig
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *ZapLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *ZapLogger) WithFields(fields Fields) ILogger {
	f := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		f = append(f, k, v)
	}
	newLogger := l.logger.With(f...)
	return &ZapLogger{
		opt:    l.opt,
		logger: newLogger,
	}
}

func (l *ZapLogger) WithError(err error) ILogger {
	newLogger := l.logger.With("error", err.Error())
	return &ZapLogger{
		opt:    l.opt,
		logger: newLogger,
	}
}

func (l *ZapLogger) Output() io.Writer {
	return io.MultiWriter(l.opt.outs...)
}

func (l *ZapLogger) callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	frames := runtime.CallersFrames([]uintptr{caller.PC})
	frame, _ := frames.Next()
	path := fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
	enc.AppendString(path)
}
