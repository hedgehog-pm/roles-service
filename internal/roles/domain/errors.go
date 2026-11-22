package domain

type Code string

const CodeInvalidInput Code = "invalid_input"

type Error struct {
	Code    Code
	Message string
}

func NewError(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}
