package dao

import (
	"database/sql"

	"github.com/jackc/pgconn"
)

const (
	_uniqueViolation = "23505"
)

func ErrUniqueViolation(err error) bool {
	return matchError(err, _uniqueViolation)
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

func matchError(err error, code string, more ...string) bool {
	if pge, ok := err.(*pgconn.PgError); ok {
		if pge.Code == code {
			return true
		}
		if n := len(more); n > 0 {
			for i := 0; i < n; i++ {
				if pge.Code == more[i] {
					return true
				}
			}
		}
	}
	return false
}
