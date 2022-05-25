package location

import (
	"encoding/json"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/crossing"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"golang.org/x/net/context"
)

func persistCrossing(ctx context.Context, c *detector.Crossing) error {
	tollgateCrossing := crossing.NewTollgateCrossing(c.TollgateID, c.WorkerID, c)
	err := tollgateCrossing.Insert(ctx)
	if err != nil {
		return fmt.Errorf("crossing persist error: %s", err)
	}
	return nil
}

func publishCrossing(ctx context.Context, c *detector.Crossing) error {
	channel := crossingChannel(c.TollgateID)

	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("crossing marshal error: %q", err)
	}
	payload := string(bytes)
	pub := dao.Reds.PubSubClient
	val, err := pub.Publish(ctx, channel, payload).Result()
	_ = val
	if err != nil {
		return fmt.Errorf("publish crossing error: %q", err)
	}
	return nil
}

func crossingChannel(tollgateID string) string {
	return fmt.Sprintf("channel:crossing:%s", tollgateID)
}
