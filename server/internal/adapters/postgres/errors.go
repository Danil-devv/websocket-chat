package postgres

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"server/internal/repository/errs"
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

func newPostgresError(e error) *Error {
	switch {
	case errors.Is(e, pgx.ErrNoRows):
		return &Error{err: errs.ErrNotFound, msg: e.Error()}
	default:
		return &Error{err: errs.ErrInternal, msg: e.Error()}
	}
}
