package worker

import (
	"testing"

	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/dom"
	"github.com/driverscooperative/geosrv/stat"
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
		wrk := newWorker(Offline)
		require.NoError(t, wrk.Insert(ctx))
	}
	for i := 0; i < 5; i++ {
		wrk := newWorker(Paused)
		require.NoError(t, wrk.Insert(ctx))
	}
	for i := 0; i < 7; i++ {
		wrk := newWorker(OnJob)
		require.NoError(t, wrk.Insert(ctx))
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
