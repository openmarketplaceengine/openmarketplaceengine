// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
)

type RawSQL struct {
	sql []string
}

//-----------------------------------------------------------------------------

func (r *RawSQL) Appendf(format string, args ...interface{}) *RawSQL {
	r.sql = append(r.sql, fmt.Sprintf(format, args...))
	return r
}

//-----------------------------------------------------------------------------

func (r *RawSQL) Exec(sql ...string) (err error) {
	if failInit(&err) {
		return
	}
	if len(sql) > 0 {
		r.sql = append(r.sql, sql...)
	}
	ctx := cfg.Context()
	db := DB()
	for i := 0; i < len(r.sql) && err == nil; i++ {
		_, err = db.ExecContext(ctx, r.sql[i])
	}
	return
}
