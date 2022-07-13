// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

type UpsertStatus int

const (
	UpsertUnknown UpsertStatus = 0
	UpsertCreated UpsertStatus = 1
	UpsertUpdated UpsertStatus = 2
)

func Upsert(ctx Context, insert, update func() Executable) (Result, UpsertStatus, error) {
	sql := insert()
	err := ExecTX(SkipErrorsContext(ctx, ErrUniqueViolation), sql)
	if err == nil && RowsAffected(sql.Result()) > 0 {
		return sql.Result(), UpsertCreated, nil
	}
	if err != nil && !ErrUniqueViolation.Is(err) {
		return nil, UpsertUnknown, err
	}
	sql = update()
	err = ExecTX(ctx, sql)
	if err == nil {
		return sql.Result(), UpsertUpdated, nil
	}
	return nil, UpsertUnknown, err
}
