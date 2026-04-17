package domain

// Error representa um erro de domínio
type Error struct {
	Message string
	Code    string
}

func (e *Error) Error() string {
	return e.Message
}

// NewError cria um novo erro de domínio
func NewError(message, code string) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}
