// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package uri

import "path"

func Join(addr string, pathElem ...string) string {
	if len(pathElem) == 0 {
		return addr
	}
	full := path.Join(pathElem...)
	plen := len(full)
	if plen == 0 || (plen == 1 && (full[0] == '.' || full[0] == '/')) {
		return addr
	}
	alen := len(addr)
	if alen == 0 {
		return full
	}
	asep := addr[alen-1] == '/'
	psep := full[0] == '/'
	if asep && psep {
		return addr + full[1:]
	}
	if asep || psep {
		return addr + full
	}
	return addr + "/" + full
}
