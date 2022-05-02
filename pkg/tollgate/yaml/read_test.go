package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTollgates(t *testing.T) {
	t.Run("testRead", func(t *testing.T) {
		testRead(t)
	})
}

func testRead(t *testing.T) {
	tg, err := ReadEmbedded()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(tg), 17)
}
