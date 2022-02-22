// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package log

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func YAML(v interface{}) string {
	var w AppendWriter
	w.Alloc(32)
	e := yaml.NewEncoder(&w)
	err := e.Encode(v)
	if err == nil {
		if err = e.Close(); err == nil {
			return w.String()
		}
	}
	return fmt.Sprintf("YAML error: %s", err)
}
