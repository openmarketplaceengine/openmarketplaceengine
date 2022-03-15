// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
)

type (
	// SUID represents opaque UUID string coming from external sources.
	// Neither format, nor length known in advance.
	SUID    = string
	Time    = dao.Time
	Context = dao.Context
)
