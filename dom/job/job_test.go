package job

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJob(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	ctx := context.Background()

	fromLat := 40.633650
	fromLon := -74.143650

	state := fmt.Sprintf("AVAILABLE-%s", uuid.NewString())
	job1 := &Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.636916,
		PickupLon: -74.195995,
	}
	job2 := &Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.634408,
		PickupLon: -74.198356,
	}

	job3 := &Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.631076,
		PickupLon: -74.167156,
	}

	job4 := &Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.632900,
		PickupLon: -74.118337,
	}

	//most faraway job
	job5 := &Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.647159,
		PickupLon: -74.320006,
	}

	//closest job
	job6 := &Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.633836,
		PickupLon: -74.142222,
	}

	for _, j := range []*Job{job1, job2, job3, job4, job5, job6} {
		_, _, err = j.Upsert(ctx)
		require.NoError(t, err)
	}

	t.Run("testQueryByPickupDistance", func(t *testing.T) {
		testQueryByPickupDistance(t, state, fromLat, fromLon)
	})

	t.Run("testJobInRangeSql", func(t *testing.T) {
		testJobInRangeSql(t)
	})
}

func testQueryByPickupDistance(t *testing.T, state string, fromLat float64, fromLon float64) {
	ctx := context.Background()

	jobs0, err := QueryByPickupDistance(ctx, fromLat, fromLon, state, 10, Km, 100)
	require.NoError(t, err)
	require.Len(t, jobs0, 5)

	jobs1, err := QueryByPickupDistance(ctx, fromLat, fromLon, state, 4.5, Mile, 100)
	require.NoError(t, err)
	require.Len(t, jobs1, 5)

	jobs2, err := QueryByPickupDistance(ctx, fromLat, fromLon, state, 20, Km, 100)
	require.NoError(t, err)
	require.Len(t, jobs2, 6)

	jobs3, err := QueryByPickupDistance(ctx, fromLat, fromLon, state, 0.5, Km, 100)
	require.NoError(t, err)
	require.Len(t, jobs3, 1)
}

func testJobInRangeSql(t *testing.T) {
	s0 := jobsInRangeSql(78.3232, 65.3234, "AVAILABLE", 4000, Mile, 20)
	require.NotContains(t, s0, "6371")
	require.Contains(t, s0, "3959")
	require.Contains(t, s0, "cos(radians(78.3232))")
	require.Contains(t, s0, "* cos(radians(pickup_lon) - radians(65.3234))")
	require.Contains(t, s0, "+ sin(radians(78.3232)) * sin(radians(pickup_lat))")
	require.Contains(t, s0, "t.distance < 4000")
	require.Contains(t, s0, "limit 20")

	s1 := jobsInRangeSql(78.3232, 65.3234, "AVAILABLE", 4000, Km, 20)
	require.NotContains(t, s1, "3959")
	require.Contains(t, s1, "6371")

}
