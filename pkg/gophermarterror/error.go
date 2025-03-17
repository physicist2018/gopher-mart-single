package gophermarterror

import "strings"

// GopherMartError - детализация ошибок сервиса GopherMart
type GopherMartError struct {
	Message string
	Err     error
}

func NewGophermartError(msg string, err error) *GopherMartError {
	return &GopherMartError{
		Message: msg,
		Err:     err,
	}
}

func (ge *GopherMartError) Error() string {
	return strings.Join([]string{ge.Message, ge.Err.Error()}, " ")
}

func (ge *GopherMartError) Unwrap() error {
	return ge.Err
}
