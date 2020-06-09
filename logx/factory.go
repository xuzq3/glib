package logx

type LoggerType int

const (
	Logrus LoggerType = iota
	Zap
)

type LoggerFactory struct {
}

func NewLoggerFactory() *LoggerFactory {
	return &LoggerFactory{}
}

func (f *LoggerFactory) Create(t LoggerType, opt *Option) ILogger {
	switch t {
	case Logrus:
		return NewLogrusLogger(opt)
	case Zap:
		return NewZapLogger(opt)
	default:
		return nil
	}
}
