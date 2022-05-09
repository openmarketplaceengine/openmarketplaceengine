package worker

import (
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/stat"
	"github.com/stretchr/testify/require"
)

func TestWorkerStat(t *testing.T) {
	dom.WillTest(t, "test", true)

	t.Run("testWorkersByStatus", func(t *testing.T) {
		testWorkersByStatus(t)
	})
}

func testWorkersByStatus(t *testing.T) {
	ctx := cfg.Context()
	for i := 0; i < 3; i++ {
		wrk := genWorker(Offline)
		require.NoError(t, wrk.Persist(ctx))
	}
	for i := 0; i < 5; i++ {
		wrk := genWorker(Paused)
		require.NoError(t, wrk.Persist(ctx))
	}
	for i := 0; i < 7; i++ {
		wrk := genWorker(OnJob)
		require.NoError(t, wrk.Persist(ctx))
	}

	status, err := workersByStatus(ctx)
	require.NoError(t, err)

	kv, ok := status.(*stat.IntKeyVal)
	require.True(t, ok)
	require.Equal(t, "offline", kv.Key[0])
	require.Equal(t, int64(3), kv.Val[0])

	require.Equal(t, "ready", kv.Key[1])
	require.Equal(t, int64(0), kv.Val[1])

	require.Equal(t, "onjob", kv.Key[2])
	require.Equal(t, int64(7), kv.Val[2])

	require.Equal(t, "paused", kv.Key[3])
	require.Equal(t, int64(5), kv.Val[3])

	require.Equal(t, "disabled", kv.Key[4])
	require.Equal(t, int64(0), kv.Val[4])

	require.Equal(t, "total", kv.Key[5])
	require.Equal(t, int64(15), kv.Val[5])
}
