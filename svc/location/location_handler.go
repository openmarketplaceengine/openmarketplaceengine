package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/worker"
)

func persistLocation(ctx context.Context, areaKey string, l *Location) error {
	stamp := dom.Time{}
	stamp.Now()
	location := worker.Location{
		Worker:    l.WorkerID,
		Stamp:     stamp,
		Longitude: l.Longitude,
		Latitude:  l.Latitude,
	}
	err := location.Persist(ctx)
	if err != nil {
		return fmt.Errorf("location persist error: %s", err)
	}
	return nil
}

func publishLocation(ctx context.Context, areaKey string, l *Location) error {
	channel := locationChannel(l.WorkerID)

	bytes, err := json.Marshal(Location{
		WorkerID:  l.WorkerID,
		Longitude: l.Longitude,
		Latitude:  l.Latitude,
	})
	if err != nil {
		return fmt.Errorf("location marshal error: %q", err)
	}
	payload := string(bytes)
	pub := dao.Reds.PubSubClient
	err = pub.Publish(ctx, channel, payload).Err()

	if err != nil {
		return fmt.Errorf("location publish error: %q", err)
	}
	return nil
}

func locationChannel(workerID string) string {
	return fmt.Sprintf("channel-location-%s", workerID)
}
