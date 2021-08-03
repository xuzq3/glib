package multicast

import "errors"

const (
	OKCode = 0
	OKMsg  = "ok"
)

var (
	ErrUnknownCmd = errors.New("unknown command")
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

var (
	ErrServerError       = NewError(10000, "server error")
	ErrInvalidConnection = NewError(10001, "invalid connection")
	ErrHandleTimeout     = NewError(10002, "handle timeout")
	ErrInvalidParams     = NewError(10003, "invalid params")
)
