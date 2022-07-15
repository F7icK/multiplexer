package infrastruct

import (
	"errors"
	"net/http"
)

type CustomError struct {
	msg  string
	Code int
}

func NewError(msg string, code int) *CustomError {
	return &CustomError{
		msg:  msg,
		Code: code,
	}
}

func (c *CustomError) Error() string {
	return c.msg
}

var (
	ErrMethodNotAllowed    = NewError("method not allowed, only POST", http.StatusMethodNotAllowed)
	ErrBadRequest          = NewError("bad query input", http.StatusBadRequest)
	ErrCountURL            = NewError("for a correct request, you need to send from 1 to 20 URLs", http.StatusRequestEntityTooLarge)
	ErrBadJSONURL          = NewError("invalid input request. json must contain the URL", http.StatusBadRequest)
	ErrHTTPLimitConnection = NewError("limit of received requests exceeded", http.StatusTooManyRequests)
	ErrValidationPort      = errors.New("bar config, port contains errors")
	ErrLimitGoRoutines     = errors.New("bar config, LimitGoRoutines cannot be less than 1")
	ErrLimitConnection     = errors.New("bar config, LimitConnection cannot be less than 1")
	ErrTimeoutIncoming     = errors.New("bar config, TimeoutIncoming cannot be less than 1")
	ErrTimeoutOutgoing     = errors.New("bar config, TimeoutOutgoing cannot be less than 1")
)
