package app

import (
	errs "chat/internal/repository/errs"
	"errors"
	"fmt"
)

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("data not found")
)

type Error struct {
	err error
	msg string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.err, e.msg)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) WithMessage(msg string) *Error {
	e.msg = fmt.Sprintf("%s; %s", e.msg, msg)
	return e
}

func newAppError(e error) *Error {
	switch {
	case errors.Is(e, errs.ErrNotFound):
		return &Error{err: ErrNotFound, msg: e.Error()}
	default:
		return &Error{err: ErrInternal, msg: e.Error()}
	}
}
