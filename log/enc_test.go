// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"testing"
	"time"
)

func TestYAML(t *testing.T) {
	m := map[string]interface{}{
		"key7": "seven",
		"key6": 6,
		"key5": time.Now(),
		"key4": []int{0, 1, 2, 3, 4},
		"key3": map[string]interface{}{"z": "zero", "a": "alpha", "b": "beta"},
	}
	fmt.Print(YAML(m))
}
