package dao

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
)

type ErrCode string

const (
	ErrUniqueViolation ErrCode = "23505"
	ErrUndefinedTable  ErrCode = "42P01"
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
