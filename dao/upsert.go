// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

type UpsertStatus int

const (
	UpsertUnknown = 0
	UpsertCreated = 1
	UpsertUpdated = 2
)

func Upsert(ctx Context, insert, update func() Executable) (Result, UpsertStatus, error) {
	sql := insert()
	err := ExecTX(ctx, sql)
	if err == nil {
		return sql.Result(), UpsertCreated, nil
	}
	if !ErrUniqueViolation(err) {
		return nil, UpsertUnknown, err
	}
	sql = update()
	err = ExecTX(ctx, sql)
	if err == nil {
		return sql.Result(), UpsertUpdated, nil
	}
	return nil, UpsertUnknown, err
}
