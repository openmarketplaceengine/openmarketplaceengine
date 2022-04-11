package htp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoutesArrayFind(t *testing.T) {
	data := map[string]int{
		"z": 3,
		"c": 2,
		"b": 1,
		"a": 0,
	}
	var ary routeArray
	for k := range data {
		require.Equal(t, -1, ary.find(k))
		ary.getOrAdd(k)
	}
	for k, v := range data {
		require.Equal(t, v, ary.find(k))
	}
}
