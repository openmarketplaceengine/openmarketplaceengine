package dao

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
)

type ErrCode string

const (
	ErrUniqueViolation ErrCode = "23505"
	ErrUndefinedTable  ErrCode = "42P01"
	ErrUndefinedColumn ErrCode = "42703"
)

//-----------------------------------------------------------------------------

func (c ErrCode) Is(err error) bool {
	for err != nil {
		if pge, ok := err.(*pgconn.PgError); ok {
			return pge.Code == string(c)
		}
		err = errors.Unwrap(err)
	}
	return false
}

//-----------------------------------------------------------------------------

func SkipNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

//-----------------------------------------------------------------------------

func WrapNoRows(err error) error {
	if err == nil {
		return sql.ErrNoRows
	}
	return err
}

//-----------------------------------------------------------------------------
// Skip Errors Context
//-----------------------------------------------------------------------------

type skipErrKey struct{}

func SkipErrorsContext(parent Context, skip ...ErrCode) Context {
	return context.WithValue(parent, skipErrKey{}, skip)
}

func SkipUndefErrors(ctx Context) Context {
	return SkipErrorsContext(ctx, ErrUndefinedColumn, ErrUndefinedTable)
}

func ShouldSkipError(ctx Context, err error) bool {
	if skip, ok := ctx.Value(skipErrKey{}).([]ErrCode); ok {
		for i := 0; i < len(skip); i++ {
			if skip[i].Is(err) {
				return true
			}
		}
	}
	return false
}
