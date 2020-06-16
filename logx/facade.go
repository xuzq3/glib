package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuzq3/glib/logx/file"
)

const (
	FileRotateByTime = "time"
	FileRotateBySize = "size"
	FileRotateByMix  = "mix"
)

type FileConfig struct {
	Enable       bool
	Filename     string
	RotateType   string // the type of rotation, support time and size
	RotateTimeS  int    // the number of seconds between rotation
	RotateSize   int    // the maximum size in megabytes of the log file before it gets rotated
	RotateCount  int    // the maximum number of log files to retain
	RotateMaxAge int    // the maximum number of days to retain
}

type Config struct {
	Level          string
	DisableConsole bool
	JsonFormat     bool
	File           FileConfig
}

var _facade *loggerFacade

func init() {
	_ = Init(Config{})
}

func Init(config Config) error {
	opt := NewOption()
	opt.AddCallerSkip(1)
	opt.SetLevel(config.Level)
	if !config.DisableConsole {
		opt.AddOutput(os.Stdout)
	}
	if config.JsonFormat {
		opt.SetJsonFormat()
	}
	if config.File.Enable {
		dir := filepath.Dir(config.File.Filename)
		err := ensureDir(dir)
		if err != nil {
			return err
		}

		switch strings.ToLower(config.File.RotateType) {
		case FileRotateBySize:
			f := file.NewSizeRotateFile(
				config.File.Filename,
				config.File.RotateSize,
				config.File.RotateCount,
				config.File.RotateMaxAge,
			)
			opt.AddOutput(f)
		case FileRotateByTime:
			f, err := file.NewTimeRotateFile(
				config.File.Filename,
				time.Second*time.Duration(config.File.RotateTimeS),
				config.File.RotateCount,
				config.File.RotateMaxAge,
			)
			if err != nil {
				return err
			}
			opt.AddOutput(f)
		case FileRotateByMix:
			f := file.NewMixRotateFile(
				config.File.Filename,
				config.File.RotateSize,
				time.Second*time.Duration(config.File.RotateTimeS),
				config.File.RotateMaxAge,
				config.File.RotateCount,
			)
			opt.AddOutput(f)
		default:
			return fmt.Errorf("unsupport file rotate type")
			//f, err := os.OpenFile(config.File.Filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			//if err != nil {
			//	return err
			//}
			//opt.AddOutput(f)
		}
	}

	_facade = &loggerFacade{
		logger: NewZapLogger(opt),
		//logger: NewLogrusLogger(opt),
	}
	return nil
}

func Debug(args ...interface{}) {
	_facade.logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	_facade.logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	_facade.logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	_facade.logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	_facade.logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	_facade.logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	_facade.logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	_facade.logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	_facade.logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	_facade.logger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	_facade.logger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	_facade.logger.Fatalf(format, args...)
}

func WithFields(fields Fields) ILogger {
	return &loggerFacade{
		logger: _facade.logger.WithFields(fields),
	}
}

func WithError(err error) ILogger {
	return &loggerFacade{
		logger: _facade.logger.WithError(err),
	}
}

func Output() io.Writer {
	return _facade.logger.Output()
}

func ensureDir(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !f.IsDir() {
			return fmt.Errorf("dir is invalid")
		}
	}
	return nil
}

type loggerFacade struct {
	logger ILogger
}

func (l *loggerFacade) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *loggerFacade) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *loggerFacade) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *loggerFacade) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *loggerFacade) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *loggerFacade) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *loggerFacade) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *loggerFacade) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *loggerFacade) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *loggerFacade) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *loggerFacade) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *loggerFacade) Panicf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *loggerFacade) WithFields(fields Fields) ILogger {
	return &loggerFacade{
		logger: l.logger.WithFields(fields),
	}
}

func (l *loggerFacade) WithError(err error) ILogger {
	return &loggerFacade{
		logger: l.logger.WithError(err),
	}
}

func (l *loggerFacade) Output() io.Writer {
	return l.logger.Output()
}
