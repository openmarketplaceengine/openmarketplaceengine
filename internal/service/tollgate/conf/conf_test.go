package conf

import (
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	dom.WillTest(t, "test", false)
	ctx := cfg.Context()
	t.Run("testLoadTollgates", func(t *testing.T) {
		testLoadTollgates(ctx, t)
	})
}

func testLoadTollgates(ctx cfg.SignalContext, t *testing.T) {
	err := LoadTollgates(ctx)
	require.NoError(t, err)
	tg, err := tollgate.QueryAll(ctx, 100)
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(tg), 17)

	for _, toll := range tg {
		assert.Less(t, toll.Created.UnixMilli(), time.Now().UnixMilli())
	}
}
