package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/driverscooperative/geosrv/dom/location"

	"github.com/driverscooperative/geosrv/dao"
)

func persistLocation(ctx context.Context, areaKey string, l *Location) error {
	loc := location.Location{
		Worker:    l.WorkerID,
		Longitude: l.Longitude,
		Latitude:  l.Latitude,
	}
	err := loc.Insert(ctx)
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
