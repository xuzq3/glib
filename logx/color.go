package logx

type Color int

const (
	NOCOLOR Color = 0
	BLACK         = 30
	RED           = 31
	GREEN         = 32
	YELLOW        = 33
	BLUE          = 34
	PURPLE        = 35
	CYAN          = 36
	GRAY          = 37
)

var levelColors = map[Level]Color{
	TRACE: GRAY,
	DEBUG: PURPLE,
	INFO:  GREEN,
	WARN:  YELLOW,
	ERROR: RED,
	FATAL: RED,
}
