package domain

type Code string

const (
	CodeInvalidInput      Code = "invalid_input"
	CodeUnknownPermission Code = "unknown_permission_code"
	CodePermissionScope   Code = "permission_scope_mismatch"
	CodeEmptyPermissions  Code = "empty_permissions"
)

type Error struct {
	Code    Code
	Message string
	Field   string
}

func NewError(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}

func FieldError(code Code, field, message string) *Error {
	return &Error{Code: code, Field: field, Message: message}
}
