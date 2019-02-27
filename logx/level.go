package logx

type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

const DefaultLevel = TRACE

var stringLevels = map[string]Level{
	"TRACE": TRACE,
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"FATAL": FATAL,
}

var levelStrings = map[Level]string{
	TRACE: "TRACE",
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

func (level Level) String() string {
	str, ok := levelStrings[level]
	if !ok {
		return "UNKNOWN"
	}
	return str
}

func (level Level) Color() Color {
	color, ok := levelColors[level]
	if !ok {
		return NOCOLOR
	}
	return color
}

// func LevelFromString(str string) (Level, error) {
// 	level, ok := stringLevels[strings.ToUpper(str)]
// 	if !ok {
// 		return nil, errors.New("unknown level string")
// 	}
// 	return level, nil
// }

// func (level Level) Valid() bool {
// 	if level < TRACE || level > CRITICAL {
// 		return false
// 	}
// 	return true
// }
