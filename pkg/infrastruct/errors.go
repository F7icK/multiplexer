package infrastruct

import "net/http"

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
	ErrorInternalServerError = NewError("internal server error", http.StatusInternalServerError)
	ErrorBadRequest          = NewError("bad query input", http.StatusBadRequest)
	ErrorCountUrl            = NewError("for a correct request, you need to send from 1 to 20 URLs", http.StatusUnprocessableEntity)
	ErrorBadJsonUrl          = NewError("invalid input request. json must contain the URL", http.StatusBadRequest)
	ErrorBadUrl              = NewError("bad query input. can't go to url", http.StatusBadRequest)
	ErrorLimitConnection     = NewError("limit of received requests exceeded", http.StatusTooManyRequests)
)
