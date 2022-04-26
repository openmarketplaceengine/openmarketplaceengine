package csv

import (
	"context"
	"fmt"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location"

	"github.com/stretchr/testify/require"
)

const csvFile = "testdata/coopdrive-gps-pings-2022.04.06.csv"

func TestImport(t *testing.T) {
	t.Skip("manual run only")
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	tracker, err := location.NewTracker(dao.Reds.StoreClient, dao.Reds.PubSubClient)
	require.NoError(t, err)
	i := NewImporter(tracker)

	crossings, err := i.Import(context.Background(), csvFile)
	require.NoError(t, err)

	for _, crossing := range crossings {
		fmt.Printf("%v\n", crossing)
	}
}
